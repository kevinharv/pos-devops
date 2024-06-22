package transaction

import (
	"database/sql"
	"log/slog"
	"net/http"
)

func StartCheckout(logger *slog.Logger, db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Freeze transaction, don't allow more items to be added
		// Apply taxes, discounts, etc.
	}
}

func CollectPayment(logger *slog.Logger, db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Take payment
		// Apply to transaction
		// If sufficient, close transaction
	}
}