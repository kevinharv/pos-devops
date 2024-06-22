package item

import (
	"database/sql"
	"log/slog"
	"net/http"

	"github.com/kevinharv/pos-devops/server/models"
	"github.com/kevinharv/pos-devops/server/utils"
)

func ArchiveItem(logger *slog.Logger, db *sql.DB) http.HandlerFunc {
	type ArchiveRequest struct {
		ItemID int `json:"itemID"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		data, err := utils.Decode[ArchiveRequest](r)
		if err != nil {
			logger.Error("Item Request - Failed to Decode", "Decoder", err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		err = models.ArchiveItem(logger, db, data.ItemID)
		if err != nil {
			logger.Error("Item Request - Failed to Archive", "DB", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}
}
