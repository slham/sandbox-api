package dao

import (
	"database/sql"
	"fmt"
	"log/slog"
	"os"

	_ "github.com/lib/pq"
)

type Dao struct {
	DB *sql.DB
}

var db Dao

func Connect() (Dao, error) {
	var err error
	connectionString := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		os.Getenv("SBDB_HOST"),
		os.Getenv("SBDB_PORT"),
		os.Getenv("SBDB_USERNAME"),
		os.Getenv("SBDB_PASSWORD"),
		os.Getenv("SBDB_DB"),
		os.Getenv("SBDB_SSL"),
	)

	db.DB, err = sql.Open("postgres", connectionString)
	if err != nil {
		slog.Error("failed to open connection.", "err", err)
		slog.Debug(connectionString)
		return db, fmt.Errorf("failed to open connection. %w", err)
	}

	if err := db.DB.Ping(); err != nil {
		slog.Error("failed to ping database", "err", err)
		return db, fmt.Errorf("failed to ping database. %w", err)
	}

	slog.Info("successfully connected to db")
	return db, nil
}

func getDB() *sql.DB {
	return db.DB
}
