package utils

import (
	"fmt"
	"log/slog"
	"os"
	"strconv"
)

type Config struct {
	ServerAddr string
	LogLevel   slog.Level
	DBUser     string
	DBPass     string
	DBHost     string
	DBPort     int
	DBDatabase string
}

func GetConfig() *Config {
	cfg := Config{}

	addrEnv := os.Getenv("SERVER_ADDR")
	if len(addrEnv) > 0 {
		cfg.ServerAddr = fmt.Sprintf(":%s", addrEnv)
	} else {
		cfg.ServerAddr = ":8080"
	}

	levelEnv := os.Getenv("LOG_LEVEL")
	if len(levelEnv) == 0 {
		cfg.LogLevel = slog.LevelInfo
	} else if levelEnv == "ERROR" {
		cfg.LogLevel = slog.LevelError
	} else if levelEnv == "WARN" {
		cfg.LogLevel = slog.LevelWarn
	} else if levelEnv == "INFO" {
		cfg.LogLevel = slog.LevelInfo
	} else if levelEnv == "DEBUG" {
		cfg.LogLevel = slog.LevelDebug
	}

	dbUser := os.Getenv("DB_USER")
	if len(dbUser) > 0 {
		cfg.DBUser = dbUser
	} else {
		cfg.DBUser = "postgres"
	}

	dbPass := os.Getenv("DB_PASS")
	if len(dbUser) > 0 {
		cfg.DBPass = dbPass
	} else {
		cfg.DBPass = "postgres"
	}

	dbHost := os.Getenv("DB_HOST")
	if len(dbHost) > 0 {
		cfg.DBHost = dbHost
	} else {
		cfg.DBHost = "localhost"
	}

	dbDatabase := os.Getenv("DB_DATABASE")
	if len(dbDatabase) > 0 {
		cfg.DBDatabase = dbDatabase
	} else {
		cfg.DBDatabase = "pos-server"
	}

	dbPort := os.Getenv("DB_PORT")
	if len(dbPort) > 0 {
		port, err := strconv.Atoi(dbPort)
		if err != nil {
			cfg.DBPort = 5432
		}
		cfg.DBPort = port
	} else {
		cfg.DBPort = 5432
	}

	return &cfg
}
