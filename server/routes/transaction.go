package routes

import (
	"log/slog"
	"net/http"

	"github.com/kevinharv/pos-devops/server/models"
	"github.com/kevinharv/pos-devops/server/utils"
)



func ExampleTransactionHandler(logger *slog.Logger) http.HandlerFunc {
	t := models.Transaction{TransactionID: 69420, POS_ID: 42069}

	return func(w http.ResponseWriter, r *http.Request) {
		err := utils.Encode(w, r, http.StatusOK, t) 
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}
}

// Create

// Update

// Retrieve

// Archive