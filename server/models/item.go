package models

import (
	"database/sql"
	"fmt"
	"log/slog"
	"strconv"
	"strings"
)

type Item struct {
	ItemID      int     `json:"itemID"`
	CategoryID  int     `json:"categoryID"`
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Price       float64 `json:"price"`
	Archived    bool    `json:"archived"`

	CreatedDate string `json:"createdDate"`
}

func GetItemByID(logger *slog.Logger, db *sql.DB, itemID int) (item *Item, err error) {
	res, err := db.Query("SELECT * FROM items WHERE itemID = $1", itemID)
	if err != nil {
		return nil, err
	}
	defer res.Close()

	hasRecord := res.Next()
	if !hasRecord {
		return nil, fmt.Errorf("no records available")
	}
	
	dbItem := Item{}
	var priceString string
	err = res.Scan(&dbItem.ItemID, &dbItem.CategoryID, &dbItem.Name, &dbItem.Description, &priceString, &dbItem.Archived, &dbItem.CreatedDate)
	if err != nil {
		return nil, err
	}

	priceString, _ = strings.CutPrefix(priceString, "$")
	dbItem.Price, err = strconv.ParseFloat(priceString, 64)
	if err != nil {
		return nil, err
	}

	res.Close()
	return &dbItem, nil
}

func CreateItem(logger *slog.Logger, db *sql.DB, categoryID int, name string, description string, price float64) error {
	itemsResult, err := db.Query("SELECT * FROM items WHERE categoryID = $1 AND itemName = $2 AND itemDescription = $3 AND price = $4", categoryID, name, description, price)
	if err != nil {
		return err
	}
	defer itemsResult.Close()

	if itemsResult.Next() {
		return fmt.Errorf("Item already exists")
	}
	
	res, err := db.Exec("INSERT INTO items (categoryID, itemName, itemDescription, price) VALUES ($1, $2, $3, $4)", categoryID, name, description, price)
	if err != nil {
		return err
	}

	rows, err := res.RowsAffected()
	if err == nil {
		logger.Debug("DB INSERT", "Rows Affected", rows)
	}

	id, err := res.LastInsertId()
	if err == nil {
		logger.Debug("DB Item Created", "InsertID", id)
	}

	return nil
}

func UpdateName(logger *slog.Logger, db *sql.DB, itemID int, name string) error {
	res, err := db.Exec("UPDATE items SET itemName = $1 WHERE itemID = $2", name, itemID)
	if err != nil {
		return err
	}

	count, err := res.RowsAffected()
	if err == nil {
		logger.Debug("DB Item Update", "Rows Affected", count)
	}

	return nil
}

func UpdateDescription(logger *slog.Logger, db *sql.DB, itemID int, description string) error {
	res, err := db.Exec("UPDATE items SET itemDescription = $1 WHERE itemID = $2", description, itemID)
	if err != nil {
		return err
	}

	count, err := res.RowsAffected()
	if err == nil {
		logger.Debug("DB Item Update", "Rows Affected", count)
	}

	return nil
}

func UpdatePrice(logger *slog.Logger, db *sql.DB, itemID int, price float64) error {
	res, err := db.Exec("UPDATE items SET price = $1 WHERE itemID = $2", price, itemID)
	if err != nil {
		return err
	}

	count, err := res.RowsAffected()
	if err == nil {
		logger.Debug("DB Item Update", "Rows Affected", count)
	}

	return nil
}

func UpdateCategory(logger *slog.Logger, db *sql.DB, itemID int, categoryID int) error {
	res, err := db.Exec("UPDATE items SET categoryID = $1 WHERE itemID = $2", categoryID, itemID)
	if err != nil {
		return err
	}

	count, err := res.RowsAffected()
	if err == nil {
		logger.Debug("DB Item Update", "Rows Affected", count)
	}

	return nil
}

func ArchiveItem(logger *slog.Logger, db *sql.DB, itemID int) error {
	res, err := db.Exec("UPDATE items SET archived = true WHERE itemID = $1", itemID)
	if err != nil {
		return err
	}

	count, err := res.RowsAffected()
	if err == nil {
		logger.Debug("DB Item Archive", "Rows Affected", count)
	}

	return nil
}
