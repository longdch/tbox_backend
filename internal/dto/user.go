package dto

import (
	"time"
)

type User struct {
	ID          int
	PhoneNumber string
	Status      int
	CreatedAt   time.Time
	UpdatedAt   time.Time
}
