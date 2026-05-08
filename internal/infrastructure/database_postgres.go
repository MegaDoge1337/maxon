package infrastructure

import (
	"database/sql"
	"log/slog"
)

func NewPostgresConnection(driver string, dsn string) (*sql.DB, error) {
	// open db connection
	db, err := sql.Open(driver, dsn)
	if err != nil {
		slog.Error("failed to open database connection", slog.Any("error", err))
		return nil, err
	}
	// test database connection
	if err := db.Ping(); err != nil {
		slog.Error("failed to ping database", slog.Any("error", err))
		return nil, err
	}
	slog.Info("database connection established")

	return db, nil
}
