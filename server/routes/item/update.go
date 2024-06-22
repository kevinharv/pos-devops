package item

import (
	"database/sql"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/kevinharv/pos-devops/server/models"
	"github.com/kevinharv/pos-devops/server/utils"
)

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

			models.UpdateCategory(logger, db, data.ItemID, int(categoryNumber))


		} else if data.Attribute == "description" {
			
			err = models.UpdateDescription(logger, db, data.ItemID, data.Value)

		} else if data.Attribute == "price" {
			priceFloat, err := strconv.ParseFloat(data.Value, 64)
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				return
			}
			
			models.UpdatePrice(logger, db, data.ItemID, priceFloat)


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
