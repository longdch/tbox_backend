package models

import (
	"tbox_backend/internal/dto"
	"time"
)

type UserOtp struct {
	UserOtpID int       `db:"user_otp_id"`
	UserID    int       `db:"user_id"`
	Otp       string    `db:"otp"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}

func (u UserOtp) ToOtp() dto.UserOtp {
	return dto.UserOtp{
		ID:        u.UserOtpID,
		UserID:    u.UserID,
		Otp:       u.Otp,
		CreatedAt: u.CreatedAt,
		UpdatedAt: u.UpdatedAt,
	}
}

func (u *UserOtp) FromDto(userOtpDto dto.UserOtp) {
	u.UserOtpID = userOtpDto.ID
	u.UserID = userOtpDto.UserID
	u.Otp = userOtpDto.Otp
	u.CreatedAt = userOtpDto.CreatedAt
	u.UpdatedAt = userOtpDto.UpdatedAt
}
