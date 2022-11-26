package models

import "time"

type Log struct {
	UserID    uint      `json:"user_id"`
	Page      string    `json:"page"`
	StartedAt time.Time `json:"started_at"`
	EndedAt   time.Time `json:"ended_at"`
}
