package transaction

import (
	"database/sql"
	"log/slog"
	"net/http"

	"github.com/kevinharv/pos-devops/server/models"
	"github.com/kevinharv/pos-devops/server/utils"
)

func RemoveItemFromTransaction(logger *slog.Logger, db *sql.DB) http.HandlerFunc {
	type RemoveItemRequest struct {
		EntryID int `json:"entryID"`
	}
	
	return func(w http.ResponseWriter, r *http.Request) {
		
		data, err := utils.Decode[RemoveItemRequest](r)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		err = models.RemoveItemFromTransaction(logger, db, data.EntryID)
		if err != nil {
			logger.Error("Transaction - Failed to remove item from transaction", "DB", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		
		// TODO - recalculate total

		w.WriteHeader(http.StatusOK)
	}
}