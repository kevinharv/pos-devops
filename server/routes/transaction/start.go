package transaction

import (
	"database/sql"
	"log/slog"
	"net/http"

	"github.com/kevinharv/pos-devops/server/models"
	"github.com/kevinharv/pos-devops/server/utils"
)

func StartTransaction(logger *slog.Logger, db *sql.DB) http.HandlerFunc {
	type StartTransactionRequest struct {
		POS_ID int `json:"posID"`
		StoreID int `json:"storeID"`
		UserID int `json:"userID"`
	}
	
	return func(w http.ResponseWriter, r *http.Request) {
		data, err := utils.Decode[StartTransactionRequest](r)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		transaction, err := models.CreateTransaction(logger, db, data.POS_ID, data.StoreID, data.UserID)
		if err != nil {
			logger.Error("Transaction - Failed to start", "DB", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		err = utils.Encode(w, r, http.StatusCreated, transaction)
		if err != nil {
			logger.Error("Transaction - Failed to encode", "Encoder", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}
}