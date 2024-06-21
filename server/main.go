/*
- Create HTTP endpoints for transactions
- Handle service-service auth
- Interact with DB
- Logging all transactions and expose metrics

Rough idea of the app flow
1. Start transaction - transID, posID, more?
2. Receive barcode - transID, posID, barcode
3. Checkout - transID, posID, more?
*/

package main

import (
	"context"
	"crypto/tls"
	"database/sql"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/kevinharv/pos-devops/server/middleware"
	"github.com/kevinharv/pos-devops/server/models"
	"github.com/kevinharv/pos-devops/server/routes"
	"github.com/kevinharv/pos-devops/server/utils"
)

func main() {
	// Get config, setup logging, DB
	config := utils.GetConfig()

	jsonHandler := slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: config.LogLevel,
	})
	logger := slog.New(jsonHandler)

	// Setup TLS
	cert, err := tls.LoadX509KeyPair("/bin/cert.pem", "/bin/key.pem")
	if err != nil {
		logger.Error("TLS Error", "TLS", err)
		logger.Error("TLS Failed - Exiting")
		os.Exit(1)
	}

	// Create a TLS Config with the certificate
	tlsConfig := &tls.Config{
		Certificates: []tls.Certificate{cert},
	}

	// Establish DB connection
	var db *sql.DB

	for {
		db, err = models.CreateConnection(config, logger)

		if err == nil {
			break
		}

		logger.Error("Failed to connect to database")
		time.Sleep(3 * time.Second)
	}

	// Configure Routes, add middleware
	mux := http.NewServeMux()
	routes.AddRoutes(mux, logger, db)
	loggedMux := middleware.LogRequest(mux, logger)
	corsMux := middleware.CORSMiddleware(loggedMux)

	s := &http.Server{
		Addr:         config.ServerAddr,
		Handler:      corsMux,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		TLSConfig:    tlsConfig,
	}

	// Create channel and handle OS signals
	exitChannel := make(chan os.Signal, 1)
	signal.Notify(exitChannel, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	// Run HTTP server in Goroutine
	go func(l *slog.Logger) {
		err := s.ListenAndServeTLS("", "") // certs loaded in config
		if err != nil && err != http.ErrServerClosed {
			l.Error("HTTP server closed unexpectedly\n")
			l.Error(fmt.Sprintf("%s\n", err))
		}
	}(logger)

	// Log ONLY when the server has started
	logger.Info("Started POS Server")

	// Shutdown the server on OS interrupts/calls
	<-exitChannel
	logger.Info("Shutting Down POS Server")

	// Create context and give 5 seconds to timeout
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer func() {
		// Close logstreams?
		// Any extra shutdown handling
		models.CloseConnection(db, logger)
		cancel()
	}()

	// Shutdown the HTTP server
	if err := s.Shutdown(ctx); err != nil {
		logger.Error("Server Shutdown Failed:%+v", err)
	}
	logger.Info("Server Shutdown Properly")
}
