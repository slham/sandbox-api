package main

import (
	"context"
	"crypto/tls"
	"log"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/slham/sandbox-api/auth"
	"github.com/slham/sandbox-api/crypt"
	"github.com/slham/sandbox-api/dao"
	"github.com/slham/sandbox-api/handler"
	"github.com/slham/sandbox-api/middlewares"
)

const (
	SERVER_READ_TIMEOUT  = 15
	SERVER_WRITE_TIMEOUT = 15
)

func main() {
	env := os.Getenv("SANDBOX_ENVIRONMENT")
	switch env {
	case "LOCAL":
		if ok := middlewares.Initialize(middlewares.DEBUG); !ok {
			log.Fatalf("failed to initialize logging")
		}
		_, err := dao.Connect()
		if err != nil {
			log.Fatalf("failed to connect to database. %s", err)
		}
		crypt.Initialize(os.Getenv("SANDBOX_AUTH_KEY"))
		slog.Info("running on local")
	default:
		slog.Info("invalid environment", "env", env)
		os.Exit(1)
	}

	r := mux.NewRouter()

	// Middlewares
	standardSessionStore := auth.NewStandardSessionStore()
	//establishSession := middlewares.Establish(standardSessionStore)
	verifySession := middlewares.Verify(standardSessionStore)
	//terminateSession := middlewares.Terminate(standardSessionStore)
	rateLimiter := middlewares.RateLimit(env)

	r.Use(middlewares.LoggingInbound)
	r.Use(rateLimiter)

	// Controllers
	authController := handler.NewAuthController(standardSessionStore)
	userController := handler.NewUserController()
	workoutController := handler.NewWorkoutController()

	// Health APIs
	r.Methods("GET").Path("/").HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello World"))
	})

	// Auth APIs
	r.Methods("GET").Path("/auth/google/login").HandlerFunc(authController.OauthGoogleLogin)
	r.Methods("GET").Path("/auth/google/callback").HandlerFunc(middlewares.Chain(authController.OauthGoogleCallback))
	r.Methods("POST").Path("/auth/login").HandlerFunc(middlewares.Chain(authController.Login))
	//r.Methods("POST").Path("/auth/logout").HandlerFunc(middlewares.Chain(authController.Logout, terminateSession))

	// User APIs
	r.Methods("POST").Path("/users").HandlerFunc(middlewares.Chain(userController.CreateUser))
	r.Methods("GET").Path("/users").HandlerFunc(middlewares.Chain(userController.GetUsers, verifySession))
	r.Methods("GET").Path("/users/{user_id}").HandlerFunc(middlewares.Chain(userController.GetUser, verifySession))
	r.Methods("PATCH").Path("/users/{user_id}").HandlerFunc(middlewares.Chain(userController.UpdateUser, verifySession))
	r.Methods("DELETE").Path("/users/{user_id}").HandlerFunc(middlewares.Chain(userController.DeleteUser, verifySession))

	// Workouts APIs
	r.Methods("POST").Path("/users/{user_id}/workouts").HandlerFunc(middlewares.Chain(workoutController.CreateWorkout, verifySession))
	r.Methods("GET").Path("/users/{user_id}/workouts").HandlerFunc(middlewares.Chain(workoutController.GetWorkouts, verifySession))
	r.Methods("GET").Path("/users/{user_id}/workouts/{workout_id}").HandlerFunc(middlewares.Chain(workoutController.GetWorkout, verifySession))
	r.Methods("PATCH").Path("/users/{user_id}/workouts/{workout_id}").HandlerFunc(middlewares.Chain(workoutController.UpdateWorkout, verifySession))
	r.Methods("DELETE").Path("/users/{user_id}/workouts/{workout_id}").HandlerFunc(middlewares.Chain(workoutController.DeleteWorkout, verifySession))

	headersOk := handlers.AllowedHeaders([]string{
		"Access-Control-Allow-Origin",
		"Access-Control-Allow-Methods",
		"Access-Control-Allow-Headers",
		"Access-Control-Allow-Credentials",
		"Access-Control-Expose-Headers",
		"Access-Control-Max-Age",
		"Access-Control-Request-Method",
		"Access-Control-Request-Headers",
		"Accept-Encoding",
		"Connection",
		"Content-Language",
		"Content-Type",
		"Origin",
		"X-Requested-With",
	})
	originsOk := handlers.AllowedOrigins([]string{"http://localhost:3000"})
	methodsOk := handlers.AllowedMethods([]string{
		http.MethodGet,
		http.MethodPost,
		http.MethodPut,
		http.MethodPatch,
		http.MethodDelete,
		http.MethodOptions,
		http.MethodHead,
	})
	allowCredentials := handlers.AllowCredentials()

	cors := handlers.CORS(
		originsOk,
		headersOk,
		methodsOk,
		allowCredentials,
	)
	handler := cors(r)

	tlsConfig := &tls.Config{
		MinVersion:               tls.VersionTLS12,
		CurvePreferences:         []tls.CurveID{tls.CurveP521, tls.CurveP384, tls.CurveP256},
		PreferServerCipherSuites: true,
		CipherSuites: []uint16{
			tls.TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384,
			tls.TLS_ECDHE_RSA_WITH_AES_256_CBC_SHA,
			tls.TLS_RSA_WITH_AES_256_GCM_SHA384,
			tls.TLS_RSA_WITH_AES_256_CBC_SHA,
		},
	}

	srv := &http.Server{
		Addr:         ":443", //TODO: YIKES
		Handler:      handler,
		ReadTimeout:  SERVER_READ_TIMEOUT * time.Second,
		WriteTimeout: SERVER_WRITE_TIMEOUT * time.Second,
		TLSConfig:    tlsConfig,
		TLSNextProto: make(map[string]func(*http.Server, *tls.Conn, http.Handler), 0),
	}

	go func() {
		cert, key := "server.crt", "server.key"
		slog.Info("starting server")
		if err := srv.ListenAndServeTLS(cert, key); err != nil {
			log.Fatalf("failed to serve. %s", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	slog.Info("stopping server")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("failed to gracefully shutdown: %s", err)
	}
	slog.Info("server gracefully stopped")
}
