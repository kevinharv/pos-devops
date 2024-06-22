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

	mux.Handle("GET /v1/transaction/{id}", http.NotFoundHandler())
	mux.Handle("POST /v1/transaction/start", http.NotFoundHandler())
	mux.Handle("PUT /v1/transaction/item/add", http.NotFoundHandler())
	mux.Handle("DELETE /v1/transaction/item/remove", http.NotFoundHandler())
	mux.Handle("POST /v1/transaction/checkout/start", http.NotFoundHandler())
	mux.Handle("POST /v1/transaction/checkout/payment", http.NotFoundHandler())
	mux.Handle("POST /v1/transaction/close", http.NotFoundHandler())

	mux.Handle("GET /v1/item/{id}", ItemByID(logger, db))
	mux.Handle("PUT /v1/item/create", CreateItem(logger, db))
	mux.Handle("PATCH /v1/item/update", UpdateItem(logger, db))
	mux.Handle("DELETE /v1/item/archive", ArchiveItem(logger, db))
}