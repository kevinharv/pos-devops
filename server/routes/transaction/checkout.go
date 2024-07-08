package transaction

import (
	"database/sql"
	"log/slog"
	"net/http"

	"github.com/kevinharv/pos-devops/server/models"
	"github.com/kevinharv/pos-devops/server/utils"
)

func StartCheckout(logger *slog.Logger, db *sql.DB) http.HandlerFunc {
	type StartCheckoutRequest struct {
		TransactionID int `json:"transactionID"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		data, err := utils.Decode[StartCheckoutRequest](r)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		// Freeze transaction, don't allow more items to be added
		err = models.UpdateStatus(logger, db, data.TransactionID, "CHECKOUT")
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		// Apply taxes, discounts, etc.

		w.WriteHeader(http.StatusNoContent)
	}
}

func CollectPayment(logger *slog.Logger, db *sql.DB) http.HandlerFunc {
	type CheckoutPaymentRequest struct {
		TransactionID int `json:"transactionID"`
		PaymentID     int `json:"paymentID"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		data, err := utils.Decode[CheckoutPaymentRequest](r)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		// Set payment method
		err = models.SetTransactionPayment(logger, db, data.TransactionID, data.PaymentID)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		// HERE YOU WOULD CHARGE THE CARD IF IT WERE REAL

		// Assume card charged and goods paid for
		err = models.UpdateStatus(logger, db, data.TransactionID, "COMPLETE")
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusNoContent)
	}
}
