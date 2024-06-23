package models

import (
	"database/sql"
	"fmt"
	"log/slog"
	"time"
)

type TransactionEntry struct {
	EntryID       int `json:"entryID"`
	TransactionID int `json:"transactionID"`
	ItemID        int `json:"itemID"`

	CreatedAt time.Time `json:"createdAt"`
}

// Get a transaction entry
func GetTransactionEntry(logger *slog.Logger, db *sql.DB, entryID int) (*TransactionEntry, error) {
	
	var entry TransactionEntry
	row := db.QueryRow("SELECT * FROM transaction_items WHERE entryID = $1", entryID)
	err := row.Scan(&entry.EntryID, &entry.TransactionID, &entry.ItemID, &entry.CreatedAt)
	if err != nil {
		return nil, fmt.Errorf("unable to retrieve transaction entry %d", entryID)
	}

	return &entry, nil
}

// Add item to transaction
func AddItemToTransaction(logger *slog.Logger, db *sql.DB, transactionID int, itemID int) (*TransactionEntry, error) {

	// Check transaction exists
	var transactionStatus string
	row := db.QueryRow("SELECT status FROM transactions WHERE transactionID = $1", transactionID)
	err := row.Scan(&transactionStatus)
	if err != nil {
		return nil, fmt.Errorf("no transaction exists for ID %d", transactionID)
	}

	// Check transaction in valid state (ready to add items)
	if transactionStatus != "ACTIVE" {
		return nil, fmt.Errorf("transaction %d is not in ACTIVE state", transactionID)
	}

	// Check item exists
	var archived bool
	row = db.QueryRow("SELECT archived FROM items WHERE itemID = $1", itemID)
	err = row.Scan(&archived)
	if err != nil {
		return nil, fmt.Errorf("item with ID %d does not exist", itemID)
	}

	if archived {
		return nil, fmt.Errorf("item with ID %d is archived", itemID)
	}

	// Add item to transaction
	var entryID int
	row = db.QueryRow("INSERT INTO transaction_items (transactionID, itemID) VALUES ($1, $2) RETURNING entryID", transactionID, itemID)
	err = row.Scan(&entryID)
	if err != nil {
		return nil, fmt.Errorf("failed to add item %d to transaction %d", itemID, transactionID)
	}

	return GetTransactionEntry(logger, db, entryID)
}

// Remove item from transaction
