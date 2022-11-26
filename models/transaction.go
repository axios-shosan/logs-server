package models

import "time"

type Transaction struct {
	Amount uint      `json:"amount"`
	Date   time.Time `json:"date"`
}
