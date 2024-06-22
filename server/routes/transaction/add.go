package transaction

import (
	"database/sql"
	"log/slog"
	"net/http"
)

func AddItemToTransaction(logger *slog.Logger, db *sql.DB) http.HandlerFunc {
	type AddItemRequest struct {
		Foo int
	}
	
	return func(w http.ResponseWriter, r *http.Request) {
		// Check item is valid
			// Exists, not archived
		// Check transaction is valid
			// Correct status

		// Add item to transaction

		// Calculate total and update

		// Return current transaction state?
	}
}