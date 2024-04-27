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
	"database/sql"
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

func addRoutes(
	mux *http.ServeMux,
	logger *slog.Logger,
	db *sql.DB,
) {
	mux.Handle("/", http.NotFoundHandler())
	mux.Handle("/foo", middleware.LogRequest(routes.FooHandler(), logger))

	mux.Handle("/v1/transaction/start", http.NotFoundHandler())
	mux.Handle("/v1/transaction/item/add", http.NotFoundHandler())
	mux.Handle("/v1/transaction/item/remove", http.NotFoundHandler())
	mux.Handle("/v1/transaction/checkout/start", http.NotFoundHandler())
	mux.Handle("/v1/transaction/checkout/payment", http.NotFoundHandler())
	mux.Handle("/v1/transaction/close", http.NotFoundHandler())
	mux.Handle("/v1/transaction/test", routes.ExampleTransactionHandler(logger))
	
	mux.Handle("/v1/item/lookup", http.NotFoundHandler())
	mux.Handle("/v1/item/create", http.NotFoundHandler())
	mux.Handle("/v1/item/update", http.NotFoundHandler())
	mux.Handle("/v1/item/archive", http.NotFoundHandler())
}

func main() {
	// Get config, setup logging, DB
	config := utils.GetConfig()
	
	jsonHandler := slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: config.LogLevel,
	})
	logger := slog.New(jsonHandler)
	
	// Establish DB connection
	db, err := models.CreateConnection(config, logger)
	if err != nil {
		panic(err)
	}

	// Configure Routes
	mux := http.NewServeMux()
	addRoutes(mux, logger, db)

	s := &http.Server{
		Addr:         config.ServerAddr,
		Handler:      mux,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	// Create channel and handle OS signals
	exitChannel := make(chan os.Signal, 1)
	signal.Notify(exitChannel, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)


	// Run HTTP server in Goroutine
	go func(l *slog.Logger) {
		err := s.ListenAndServe()
		if err != nil && err != http.ErrServerClosed {
			l.Error("HTTP server closed unexpectedly\n")
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
