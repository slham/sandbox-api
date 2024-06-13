package main

import (
	"log"
	"log/slog"
	"os"

	"github.com/slham/sandbox-api/dao"
	"github.com/slham/toolbelt/l"
)

func main() {
	env := os.Getenv("SANDBOX_ENV")
	switch env {
	case "LOCAL":
		if ok := l.Initialize(l.DEBUG); !ok {
			log.Panicf("failed to initialize logging")
		}
		_, err := dao.Connect()
		if err != nil {
			log.Panicf("failed to connect to database. %s", err)
		}
		slog.Info("running on local")
	default:
		slog.Info("invalid environment", "env", env)
		os.Exit(1)
	}
}
