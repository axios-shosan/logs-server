package models

import "time"

type Log struct {
	UserID    uint      `json:"user_id"`
	Page      string    `json:"page"`
	StartedAt time.Time `json:"started_at"`
	EndedAt   time.Time `json:"ended_at"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
