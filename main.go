package main

import (
	"log"
	"log/slog"
	"os"

	"github.com/slham/toolbelt/l"
)

func main() {
	env := os.Getenv("SANDBOX_ENV")
	switch env {
	case "LOCAL":
		if ok := l.Initialize(l.DEBUG); !ok {
			log.Panicf("failed to initialize logging")
		}
		slog.Info("running on local")
	default:
		slog.Info("invalid environment: %s", env)
		os.Exit(1)
	}
}
