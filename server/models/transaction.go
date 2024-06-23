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

	Entries 	*[]TransactionEntry	`json:"entries"`

	StartTime   time.Time  `json:"startTime"`
	EndTime     *time.Time `json:"endTime"`
	CreatedDate time.Time  `json:"createdDate"`
}

func CreateTransaction(logger *slog.Logger, db *sql.DB, posID int, storeID int, userID int) (Transaction, error) {

	// Do on route?
	// TODO - check POS has no active or pending transactions
	// TODO - if user, check they have no active or pending transactions

	transactionStartTimestamp := time.Now()

	var transactionID int
	err := db.QueryRow("INSERT INTO transactions (status, posID, storeID, userID, startTime) VALUES ($1, $2, $3, $4, $5) RETURNING transactionID", "ACTIVE", posID, storeID, userID, transactionStartTimestamp).Scan(&transactionID)
	if err != nil {
		logger.Error("Failed to create transaction", "DB", err)
	}

	return GetTransactionByID(logger, db, transactionID)
}

func GetTransactionByID(logger *slog.Logger, db *sql.DB, transactionID int) (Transaction, error) {
	res, err := db.Query("SELECT * FROM transactions WHERE transactionID = $1", transactionID)
	if err != nil {
		return Transaction{}, err
	}

	hasRecord := res.Next()
	if !hasRecord {
		return Transaction{}, fmt.Errorf("no records available")
	}

	dbItem := Transaction{}
	var totalString string
	err = res.Scan(&dbItem.TransactionID, &dbItem.Status, &dbItem.POS_ID, &dbItem.StoreID, &dbItem.UserID, &totalString, &dbItem.PaymentID, &dbItem.Archived, &dbItem.StartTime, &dbItem.EndTime, &dbItem.CreatedDate)
	if err != nil {
		return Transaction{}, err
	}

	totalString, _ = strings.CutPrefix(totalString, "$")
	dbItem.Total, err = strconv.ParseFloat(totalString, 64)
	if err != nil {
		return Transaction{}, err
	}

	entries, err := AllTransactionEntries(logger, db, transactionID)
	if err != nil {
		logger.Error("Transaction - Failed to retrieve entries")
		return Transaction{}, err
	}

	dbItem.Entries = entries

	res.Close()
	return dbItem, nil
}
