package dto

import (
	"time"
)

type UserOtp struct {
	ID        int
	UserID    int
	Otp       string
	CreatedAt time.Time
	UpdatedAt time.Time
}
