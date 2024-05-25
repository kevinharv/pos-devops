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
		if r.Method != http.MethodPost {
			w.WriteHeader(http.StatusBadRequest)
		}

		// Decode request
		data, err := utils.Decode[UpdateRequest](r)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		// Update item name in DB
		err = models.UpdateName(logger, db, data.ItemID, data.Name)
		if err != nil {
			logger.Error("Item Request - Failed to Update", "DB", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		// Get DB item for response
		item, err := models.GetItemByID(logger, db, data.ItemID)
		if err != nil {
			logger.Error("Item Request - Failed to Validate Update", "DB", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		// Send newly updated item as response
		err = utils.Encode(w, r, http.StatusOK, item)
		if err != nil {
			logger.Error("Failed to encode response")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}
}
