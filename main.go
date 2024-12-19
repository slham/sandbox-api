package main

import (
	"context"
	"log"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gorilla/mux"
	"github.com/slham/sandbox-api/auth"
	"github.com/slham/sandbox-api/crypt"
	"github.com/slham/sandbox-api/dao"
	"github.com/slham/sandbox-api/handler"
	"github.com/slham/sandbox-api/middlewares"
	"github.com/slham/toolbelt/l"
)

const (
	SERVER_READ_TIMEOUT  = 15
	SERVER_WRITE_TIMEOUT = 15
)

func main() {
	env := os.Getenv("SANDBOX_ENVIRONMENT")
	switch env {
	case "LOCAL":
		if ok := l.Initialize(l.DEBUG); !ok {
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

	// Controllers
	authController := handler.NewAuthController()
	userController := handler.NewUserController()
	workoutController := handler.NewWorkoutController()

	r := mux.NewRouter()

	// Middlewares
	standardSessionStore := auth.NewStandardSessionStore()
	establishSession := middlewares.Establish(standardSessionStore)
	verifySession := middlewares.Verify(standardSessionStore)
	//terminateSession := middlewares.Terminate(standardSessionStore)
	rateLimiter := middlewares.RateLimit(env)

	r.Use(l.Logging)
	r.Use(rateLimiter)

	// Auth APIs
	r.Methods("GET").Path("/auth/google/login").HandlerFunc(authController.OauthGoogleLogin)
	r.Methods("GET").Path("/auth/google/callback").HandlerFunc(middlewares.Chain(authController.OauthGoogleCallback, establishSession))
	r.Methods("POST").Path("/login").HandlerFunc(middlewares.Chain(authController.Login, establishSession))
	//r.Methods("POST").Path("/logout").HandlerFunc(middlewares.Chain(authController.Logout, terminateSession))

	// User APIs
	r.Methods("POST").Path("/users").HandlerFunc(middlewares.Chain(userController.CreateUser, establishSession)) //TODO: should this be `/register`?
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

	srv := &http.Server{
		Addr:         ":8080", //TODO: YIKES
		Handler:      r,
		ReadTimeout:  SERVER_READ_TIMEOUT * time.Second,
		WriteTimeout: SERVER_WRITE_TIMEOUT * time.Second,
	}

	go func() {
		slog.Info("starting server")
		if err := srv.ListenAndServe(); err != nil {
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
