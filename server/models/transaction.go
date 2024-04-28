package models

import ()

type Transaction struct {
	TransactionID int  `json:"transactionID"`
	POS_ID        int  `json:"posID"`
	StoreID       int  `json:"storeID"`
	UserID        int  `json:"userID"`
	Total         int  `json:"total"`
	PaymentID     int  `json:"paymentID"`
	Archived      bool `json:"archived"`

	StartTime   string `json:"startTime"`
	EndTime     string `json:"endTime"`
	CreatedDate string `json:"createdDate"`
}

// CRU ops - archive in lieu of delete
