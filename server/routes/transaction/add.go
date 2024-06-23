package transaction

import (
	"database/sql"
	"log/slog"
	"net/http"

	"github.com/kevinharv/pos-devops/server/models"
	"github.com/kevinharv/pos-devops/server/utils"
)

func AddItemToTransaction(logger *slog.Logger, db *sql.DB) http.HandlerFunc {
	type AddItemRequest struct {
		TransactionID int `json:"transactionID"`
		ItemID 		  int `json:"itemID"`
	}
	
	return func(w http.ResponseWriter, r *http.Request) {
		data, err := utils.Decode[AddItemRequest](r)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		transaction, err := models.AddItemToTransaction(logger, db, data.TransactionID, data.ItemID)
		if err != nil {
			logger.Error("Transaction - Failed to add item", "DB", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		
		// TODO - Calculate total and update

		err = utils.Encode(w, r, http.StatusCreated, transaction)
		if err != nil {
			logger.Error("Transaction - Failed to encode", "Encoder", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}
}