package models

import (
	"database/sql"
	"fmt"
	"log/slog"
	"strconv"
	"strings"
)

type Transaction struct {
	TransactionID int     `json:"transactionID"`
	POS_ID        int     `json:"posID"`
	StoreID       int     `json:"storeID"`
	UserID        int     `json:"userID"`
	Total         float64 `json:"total"`
	PaymentID     int     `json:"paymentID"`
	Archived      bool    `json:"archived"`

	StartTime   string `json:"startTime"`
	EndTime     string `json:"endTime"`
	CreatedDate string `json:"createdDate"`
}

// CRU ops - archive in lieu of delete
func GetTransactionByID(logger *slog.Logger, db *sql.DB, transactionID int) (transaction Transaction, err error) {
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
	err = res.Scan(&dbItem.TransactionID, &dbItem.POS_ID, &dbItem.StoreID, &dbItem.UserID, &totalString, &dbItem.PaymentID, &dbItem.Archived, &dbItem.StartTime, &dbItem.EndTime, &dbItem.CreatedDate)
	if err != nil {
		return Transaction{}, err
	}

	totalString, _ = strings.CutPrefix(totalString, "$")
	dbItem.Total, err = strconv.ParseFloat(totalString, 64)
	if err != nil {
		return Transaction{}, err
	}

	res.Close()
	return dbItem, nil
}
