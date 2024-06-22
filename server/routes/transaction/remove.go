package transaction

import (
	"database/sql"
	"log/slog"
	"net/http"
)

func RemoveItemFromTransaction(logger *slog.Logger, db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Check item in transaction

		// Delete top 1

		// Recalculate and update total
	}
}