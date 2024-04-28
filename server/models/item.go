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

func GetItemByID(logger *slog.Logger, db *sql.DB, itemID int) (item Item, err error) {
	res, err := db.Query("SELECT * FROM items WHERE itemID = $1", itemID)
	if err != nil {
		return Item{}, err
	}

	hasRecord := res.Next()
	if !hasRecord {
		return Item{}, fmt.Errorf("no records available")
	}
	
	dbItem := Item{}
	var priceString string
	err = res.Scan(&dbItem.ItemID, &dbItem.CategoryID, &dbItem.Name, &dbItem.Description, &priceString, &dbItem.Archived, &dbItem.CreatedDate)
	if err != nil {
		return Item{}, err
	}

	priceString, _ = strings.CutPrefix(priceString, "$")
	dbItem.Price, err = strconv.ParseFloat(priceString, 64)
	if err != nil {
		return Item{}, err
	}

	res.Close()
	return dbItem, nil
}

func CreateItem(logger *slog.Logger, db *sql.DB, categoryID int, name string, description string, price float64) error {
	res, err := db.Exec("INSERT INTO items (categoryID, itemName, itemDescription, price) VALUES ($1, $2, $3, $4)", categoryID, name, description, price)
	if err != nil {
		return err
	}

	id, err := res.LastInsertId()
	if err == nil {
		logger.Debug("DB Item Created", "InsertID", id)
	}

	return nil
}

func UpdateName(logger *slog.Logger, db *sql.DB, itemID int, name string) error {
	return nil
}

func UpdateDescription(logger *slog.Logger, db *sql.DB, itemID int, description string) error {
	return nil
}

func UpdatePrice(logger *slog.Logger, db *sql.DB, itemID int, price float64) error {
	return nil
}

func UpdateCategory(logger *slog.Logger, db *sql.DB, itemID int, categoryID int) error {
	return nil
}

func ArchiveItem(logger *slog.Logger, db *sql.DB, itemID int) error {
	return nil
}
