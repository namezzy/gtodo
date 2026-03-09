package model

import (
	"time"
)

type Status string

const (
	Todo Status = "todo"
	Done Status = "done"
)

type Task struct {
	ID          int       `json:"id"`
	Description string    `json:"description"`
	Priority    string    `json:"priority"` // "high", "medium", "low"
	CreatedAt   time.Time `json:"created_at"`
	DoneAt      time.Time `json:"done_at,omitempty"`
	Status      Status    `json:"status"`
}
