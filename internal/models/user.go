package models

import (
	"tbox_backend/internal/dto"
	"time"
)

type User struct {
	UserID      int       `db:"user_id"`
	PhoneNumber string    `db:"phone_number"`
	Status      int       `db:"status"`
	CreatedAt   time.Time `db:"created_at"`
	UpdatedAt   time.Time `db:"updated_at"`
}

func (u User) ToDto() dto.User {
	return dto.User{
		ID:          u.UserID,
		PhoneNumber: u.PhoneNumber,
		Status:      u.Status,
		CreatedAt:   u.CreatedAt,
		UpdatedAt:   u.UpdatedAt,
	}
}

func (u *User) FromDto(userDto *dto.User) {
	u.UserID = userDto.ID
	u.PhoneNumber = userDto.PhoneNumber
	u.Status = userDto.Status
	u.CreatedAt = userDto.CreatedAt
	u.UpdatedAt = userDto.UpdatedAt
}
