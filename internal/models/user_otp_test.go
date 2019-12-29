package models_test

import (
	"tbox_backend/internal/dto"
	"tbox_backend/internal/models"
	"testing"
	"time"
)

func TestUserOtp_ToDto(t *testing.T) {
	now := time.Now()

	userOtpModel := models.UserOtp{
		UserOtpID: 1,
		UserID:    2,
		Otp:       "12345",
		CreatedAt: now,
		UpdatedAt: now,
	}

	expectedUserOtpDto := dto.UserOtp{
		ID:        1,
		UserID:    2,
		Otp:       "12345",
		CreatedAt: now,
		UpdatedAt: now,
	}

	if userOtpModel.UserOtpID != expectedUserOtpDto.ID ||
		userOtpModel.UserID != expectedUserOtpDto.UserID ||
		userOtpModel.Otp != expectedUserOtpDto.Otp ||
		userOtpModel.CreatedAt != expectedUserOtpDto.CreatedAt ||
		userOtpModel.UpdatedAt != expectedUserOtpDto.UpdatedAt {
		t.Fatalf("Expected: %v", expectedUserOtpDto)
	}
}

func TestUserOtp_FromDto(t *testing.T) {
	now := time.Now()

	userOtpDto := dto.UserOtp{
		ID:        1,
		UserID:    2,
		Otp:       "123456",
		CreatedAt: now,
		UpdatedAt: now,
	}

	expectedUserOtpModel := models.UserOtp{
		UserOtpID: 1,
		UserID:    2,
		Otp:       "123456",
		CreatedAt: now,
		UpdatedAt: now,
	}

	userOtpModel := &models.UserOtp{}
	userOtpModel.FromDto(userOtpDto)

	if userOtpModel.UserOtpID != expectedUserOtpModel.UserOtpID ||
		userOtpModel.UserID != expectedUserOtpModel.UserID ||
		userOtpModel.Otp != expectedUserOtpModel.Otp ||
		userOtpModel.CreatedAt != expectedUserOtpModel.CreatedAt ||
		userOtpModel.UpdatedAt != expectedUserOtpModel.UpdatedAt {
		t.Fatalf("Expected: %v", expectedUserOtpModel)
	}
}
