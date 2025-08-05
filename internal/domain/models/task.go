package models

import (
	"time"
)

type Task struct {
	Id          uint64
	Title       string
	Description string
	Priority    string
	Status      string
	CreatedAt   time.Time
	DueDate     time.Time
}
