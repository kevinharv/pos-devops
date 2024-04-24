package utils

import (
	"log/slog"
	"os"
	"fmt"
)

type Config struct {
	ServerAddr 	string
	LogLevel	slog.Level
}

func GetConfig() *Config {
	cfg := Config{}
	
	addrEnv := os.Getenv("SERVER_ADDR")
	if len(addrEnv) == 0 {
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

	return &cfg
}