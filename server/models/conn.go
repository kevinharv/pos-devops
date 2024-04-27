package models

import (
	"database/sql"
	"fmt"
	"log/slog"

	"github.com/kevinharv/pos-devops/server/utils"
	_ "github.com/lib/pq"
)

func CreateConnection(c *utils.Config, logger *slog.Logger) (*sql.DB, error) {
	connStr := fmt.Sprintf("postgresql://%s:%s@%s:%d/%s?sslmode=disable", c.DBUser, c.DBPass, c.DBHost, c.DBPort, c.DBDatabase)

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		logger.Error("Failed to open connection to database")
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		logger.Error("Failed to ping database")
		return nil, err
	}

	logger.Info("Connected to database")
	return db, nil
}

func CloseConnection(db *sql.DB, logger *slog.Logger) error {
	err := db.Close()
	if err != nil {
		logger.Error("An error occurred while closing the database connection")
		return err
	}

	return nil
}
