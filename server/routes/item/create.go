package item

import (
	"database/sql"
	"log/slog"
	"net/http"

	"github.com/kevinharv/pos-devops/server/models"
	"github.com/kevinharv/pos-devops/server/utils"
)

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

		// Maybe do price validation here? f64 <> PSQL Money might be weird

		err = models.CreateItem(logger, db, data.CategoryID, data.Name, data.Description, data.Price)
		if err != nil {
			logger.Error("Item Request - Failed to Create", "DB", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusCreated)
	}
}