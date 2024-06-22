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
		}

		err = utils.Encode(w, r, http.StatusOK, item)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}
}

func CreateItem(logger *slog.Logger, db *sql.DB) http.HandlerFunc {
	type CreateRequest struct {
		CategoryID  int     `json:"categoryID"`
		Name        string  `json:"name"`
		Description string  `json:"description"`
		Price       float64 `json:"price"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		data, err := utils.Decode[CreateRequest](r)
		if err != nil {
			logger.Error("Item Request - Failed to Decode", "Decoder", err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		// Perform validation here (price is reasonable, values are present, etc.)

		err = models.CreateItem(logger, db, data.CategoryID, data.Name, data.Description, data.Price)
		if err != nil {
			logger.Error("Item Request - Failed to Create", "DB", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusCreated)
	}
}

func UpdateItem(logger *slog.Logger, db *sql.DB) http.HandlerFunc {
	type UpdateRequest struct {
		ItemID    int    `json:"itemID"`
		Attribute string `json:"attribute"`
		Value     string `json:"value"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		// Decode request
		data, err := utils.Decode[UpdateRequest](r)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		if data.Attribute == "name" {
			err = models.UpdateName(logger, db, data.ItemID, data.Value)
		} else if data.Attribute == "category" {
			categoryNumber, err := strconv.ParseInt(data.Value, 10, 32)
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				return
			}
			err = models.UpdateCategory(logger, db, data.ItemID, int(categoryNumber))
		} else if data.Attribute == "description" {
			err = models.UpdateDescription(logger, db, data.ItemID, data.Value)
		} else if data.Attribute == "price" {
			priceFloat, err := strconv.ParseFloat(data.Value, 64)
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				return
			}
			err = models.UpdatePrice(logger, db, data.ItemID, priceFloat)
		} else {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		// Update item attribute in DB
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
