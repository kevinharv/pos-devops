package item

import (
	"database/sql"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/kevinharv/pos-devops/server/models"
	"github.com/kevinharv/pos-devops/server/utils"
)

func ItemByID(logger *slog.Logger, db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			w.WriteHeader(http.StatusBadRequest)
		}

		id, err := strconv.Atoi(r.PathValue("id"))
		if err != nil {
			logger.Error("Item Request - Invalid ID", "ID", r.PathValue("id"))
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		item, err := models.GetItemByID(logger, db, id)
		if err != nil {
			logger.Error("Item Request - Failed to Retrieve", "DB", err)
		}

		err = utils.Encode(w, r, http.StatusOK, item)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}
}