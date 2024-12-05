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
	"github.com/slham/sandbox-api/crypt"
	"github.com/slham/sandbox-api/dao"
	"github.com/slham/sandbox-api/handler"
	"github.com/slham/toolbelt/l"
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
	// r.Use(l.Logging)

	// Auth APIs
	r.Methods("GET").Path("/auth/google/login").HandlerFunc(authController.OauthGoogleLogin)
	r.Methods("GET").Path("/auth/google/callback").HandlerFunc(authController.OauthGoogleCallback)
	//r.Methods("POST").Path("/login").HandlerFunc()
	//r.Methods("POST").Path("/logout").HandlerFunc()

	// User APIs
	r.Methods("POST").Path("/users").HandlerFunc(userController.CreateUser)
	r.Methods("GET").Path("/users").HandlerFunc(userController.GetUsers)
	r.Methods("GET").Path("/users/{user_id}").HandlerFunc(userController.GetUser)
	r.Methods("PATCH").Path("/users/{user_id}").HandlerFunc(userController.UpdateUser)
	r.Methods("DELETE").Path("/users/{user_id}").HandlerFunc(userController.DeleteUser)

	// Workouts APIs
	r.Methods("POST").Path("/users/{user_id}/workouts").HandlerFunc(workoutController.CreateWorkout)
	r.Methods("GET").Path("/users/{user_id}/workouts").HandlerFunc(workoutController.GetWorkouts)
	r.Methods("GET").Path("/users/{user_id}/workouts/{workout_id}").HandlerFunc(workoutController.GetWorkout)
	r.Methods("PATCH").Path("/users/{user_id}/workouts/{workout_id}").HandlerFunc(workoutController.UpdateWorkout)
	r.Methods("DELETE").Path("/users/{user_id}/workouts/{workout_id}").HandlerFunc(workoutController.DeleteWorkout)

	srv := &http.Server{
		Addr:         ":8080",
		Handler:      r,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
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
