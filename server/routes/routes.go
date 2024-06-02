package routes

import (
	"database/sql"
	"log/slog"
	"net/http"
)

func AddRoutes(
	mux *http.ServeMux,
	logger *slog.Logger,
	db *sql.DB,
) {
	mux.Handle("/", http.NotFoundHandler())
	mux.Handle("/healthz", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	mux.Handle("/foo", FooHandler())

	mux.Handle("/v1/transaction/start", http.NotFoundHandler())
	mux.Handle("/v1/transaction/item/add", http.NotFoundHandler())
	mux.Handle("/v1/transaction/item/remove", http.NotFoundHandler())
	mux.Handle("/v1/transaction/checkout/start", http.NotFoundHandler())
	mux.Handle("/v1/transaction/checkout/payment", http.NotFoundHandler())
	mux.Handle("/v1/transaction/close", http.NotFoundHandler())
	mux.Handle("/v1/transaction/test", ExampleTransactionHandler(logger))

	mux.Handle("/v1/item/{id}", ItemByID(logger, db))
	mux.Handle("/v1/item/create", http.NotFoundHandler())
	mux.Handle("/v1/item/update", UpdateItemName(logger, db))
	mux.Handle("/v1/item/archive", http.NotFoundHandler())
}