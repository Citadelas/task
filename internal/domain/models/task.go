package models

import (
	"time"
)

type Task struct {
	Id          uint64
	Title       string
	Description string
	Status      string
	Priority    string
	CreatedAt   time.Time
	DueDate     time.Time
}
