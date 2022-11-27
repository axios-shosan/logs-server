package models

import "time"

type Transaction struct {
	Amount    uint      `json:"amount"`
	Date      time.Time `json:"date"`
	CreatedAt time.Time `json:"created_at"`
}
