package routes

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

		id, err := strconv.Atoi(r.PathValue("id"))
		if err != nil {
			logger.Error("Item Request - Invalid ID", "ID", r.PathValue("id"))
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		item, err := models.GetItemByID(logger, db, id)
		if err != nil {
			logger.Error("Item Request - Failed to Retrieve", "DB", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		err = utils.Encode(w, r, http.StatusOK, item)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}
}

func UpdateItemName(logger *slog.Logger, db *sql.DB) http.HandlerFunc {
	type UpdateRequest struct {
		ItemID int    `json:"itemID"`
		Name   string `json:"name"`
	}

	return func(w http.ResponseWriter, r *http.Request) {

		data, err := utils.Decode[UpdateRequest](r)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		err = models.UpdateName(logger, db, data.ItemID, data.Name)
		if err != nil {
			logger.Error("Item Request - Failed to Retrieve", "DB", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
	}
}
