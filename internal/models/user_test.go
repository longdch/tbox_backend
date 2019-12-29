package models_test

import (
	"tbox_backend/internal/dto"
	"tbox_backend/internal/models"
	"testing"
	"time"
)

func TestUser_ToDto(t *testing.T) {
	now := time.Now()

	userModel := models.User{
		UserID:      1,
		PhoneNumber: "0967212212",
		Status:      0,
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	expectedUserDto := dto.User{
		ID:          1,
		PhoneNumber: "0967212212",
		Status:      0,
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	if userModel.UserID != expectedUserDto.ID ||
		userModel.PhoneNumber != expectedUserDto.PhoneNumber ||
		userModel.Status != expectedUserDto.Status ||
		userModel.CreatedAt != expectedUserDto.CreatedAt ||
		userModel.UpdatedAt != expectedUserDto.UpdatedAt {
		t.Fatalf("Expected: %v", expectedUserDto)
	}
}

func TestUser_FromDto(t *testing.T) {
	now := time.Now()

	userDto := dto.User{
		ID:          1,
		PhoneNumber: "0967212212",
		Status:      0,
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	expectedUserModel := models.User{
		UserID:      1,
		PhoneNumber: "0967212212",
		Status:      0,
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	userModel := &models.User{}
	userModel.FromDto(&userDto)

	if userModel.UserID != expectedUserModel.UserID ||
		userModel.PhoneNumber != expectedUserModel.PhoneNumber ||
		userModel.Status != expectedUserModel.Status ||
		userModel.CreatedAt != expectedUserModel.CreatedAt ||
		userModel.UpdatedAt != expectedUserModel.UpdatedAt {
		t.Fatalf("Expected: %v", expectedUserModel)
	}
}
