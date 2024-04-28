package models

import ()

type Item struct {
	ItemID      int    `json:"itemID"`
	CategoryID  int    `json:"categoryID"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Price       int    `json:"price"`

	CreatedDate string `json:"createdDate"`
}

// CRUA operations for items - place behind admin middleware