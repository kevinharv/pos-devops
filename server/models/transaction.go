package models

import (
	"database/sql"
	"fmt"
	"log/slog"
	"strconv"
	"strings"
	"time"
)

type Transaction struct {
	TransactionID int     `json:"transactionID"`
	Status        string  `json:"status"`
	POS_ID        int     `json:"posID"`
	StoreID       *int    `json:"storeID"`
	UserID        *int    `json:"userID"`
	Total         float64 `json:"total"`
	PaymentID     *int    `json:"paymentID"`
	Archived      bool    `json:"archived"`

	Entries *[]TransactionEntry `json:"entries"`

	StartTime   time.Time  `json:"startTime"`
	EndTime     *time.Time `json:"endTime"`
	CreatedDate time.Time  `json:"createdDate"`
}

/*
	TRANSACTION STATUSES
	- ACTIVE
	- CHECKOUT
	- ERROR_INCOMPLETE
	- ERROR_COMPLETE
	- COMPLETE
	- CANCELLED
*/

func CreateTransaction(logger *slog.Logger, db *sql.DB, posID int, storeID int, userID int) (*Transaction, error) {

	// Ensure POS and User do not have ongoing transactions
	rows, err := db.Query("SELECT transactionID FROM transactions WHERE (posID = $1 OR userID = $2) AND status IN ('ACTIVE', 'CHECKOUT', 'ERROR_INCOMPLETE')", posID, userID)
	if err != nil {
		logger.Error("Pending transaction lookup failed", "DB", err)
		return nil, fmt.Errorf("failed to check active transactions")
	}
	defer rows.Close()

	if rows.Next() {
		logger.Info("Transaction in progress - cannot start new")
		return nil, fmt.Errorf("existing transaction in progress")
	}

	// Create the transaction
	transactionStartTimestamp := time.Now()
	var transactionID int
	err = db.QueryRow("INSERT INTO transactions (status, posID, storeID, userID, startTime) VALUES ($1, $2, $3, $4, $5) RETURNING transactionID", "ACTIVE", posID, storeID, userID, transactionStartTimestamp).Scan(&transactionID)
	if err != nil {
		logger.Error("Failed to create transaction", "DB", err)
	}

	// Return the newly created transaction
	return GetTransactionByID(logger, db, transactionID)
}

func GetTransactionByID(logger *slog.Logger, db *sql.DB, transactionID int) (*Transaction, error) {
	res, err := db.Query("SELECT * FROM transactions WHERE transactionID = $1", transactionID)
	if err != nil {
		return nil, err
	}
	defer res.Close()

	hasRecord := res.Next()
	if !hasRecord {
		return nil, fmt.Errorf("no records available")
	}

	dbItem := Transaction{}
	var totalString string
	err = res.Scan(&dbItem.TransactionID, &dbItem.Status, &dbItem.POS_ID, &dbItem.StoreID, &dbItem.UserID, &totalString, &dbItem.PaymentID, &dbItem.Archived, &dbItem.StartTime, &dbItem.EndTime, &dbItem.CreatedDate)
	if err != nil {
		return nil, err
	}

	totalString, _ = strings.CutPrefix(totalString, "$")
	dbItem.Total, err = strconv.ParseFloat(totalString, 64)
	if err != nil {
		return nil, err
	}

	entries, err := AllTransactionEntries(logger, db, transactionID)
	if err != nil {
		logger.Error("Transaction - Failed to retrieve entries")
		return nil, err
	}

	dbItem.Entries = entries
	return &dbItem, nil
}

func UpdateStatus(logger *slog.Logger, db *sql.DB, transactionID int, status string) error {
	if !validateStatus(status) {
		return fmt.Errorf("invalid status for %d", transactionID)
	}

	_, err := db.Exec("UPDATE transactions SET status = $1 WHERE transactionID = $2", status, transactionID)
	if err != nil {
		return fmt.Errorf("error when updating status: %s", err)
	}

	return nil
}

func validateStatus(status string) bool {
	switch status {
	case "ACTIVE", "CHECKOUT", "ERROR_INCOMPLETE", "ERROR_COMPLETE", "COMPLETE", "CANCELLED":
		return true
	}
	return false
}
