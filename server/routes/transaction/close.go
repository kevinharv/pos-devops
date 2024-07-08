package transaction

import (
	"database/sql"
	"log/slog"
	"net/http"

	"github.com/kevinharv/pos-devops/server/models"
	"github.com/kevinharv/pos-devops/server/utils"
)

func CloseTransaction(logger *slog.Logger, db *sql.DB) http.HandlerFunc {
	type CloseTransactionRequest struct {
		TransactionID int `json:"transactionID"`
	}
	
	return func(w http.ResponseWriter, r *http.Request) {
		data, err := utils.Decode[CloseTransactionRequest](r)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		transaction, err := models.GetTransactionByID(logger, db, data.TransactionID)
		if err != nil {
			w.WriteHeader(http.StatusNotFound)
			return
		}

		if transaction.Status == "ACTIVE" || transaction.Status == "CHECKOUT" || transaction.Status == "CLOSED" {
			// Cancel transaction
			err = models.UpdateStatus(logger, db, data.TransactionID, "CANCELLED")
			if err != nil {
				logger.Error("Trasaction Cancel Failed", "DB", err)
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
		} else {
			// Move to error state
			err = models.UpdateStatus(logger, db, data.TransactionID, "ERROR_INCOMPLETE")
			if err != nil {
				logger.Error("Trasaction Close Failed", "DB", err)
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
		}
	}
}