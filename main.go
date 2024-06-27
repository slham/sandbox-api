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
	"github.com/slham/sandbox-api/dao"
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
		slog.Info("running on local")
	default:
		slog.Info("invalid environment", "env", env)
		os.Exit(1)
	}

    r := mux.NewRouter()

    //TODO - add api routes

    srv := &http.Server{
        Addr:":8080",
        Handler: r,
        ReadTimeout: 15 * time.Second,
        WriteTimeout: 15 * time.Second,
    }

    go func(){
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
