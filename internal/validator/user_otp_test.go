package validator_test

import (
	"tbox_backend/internal/validator"
	"testing"
)

func TestUserOtpValidator_IsOtpValid_Success(t *testing.T) {
	var userOtpValidator validator.IUserOtpValidator
	userOtpValidator = validator.UserOtpValidator{}

	if !userOtpValidator.IsOtpValid("123456", 6) {
		t.Fatal("expected true")
	}
}

func TestUserOtpValidator_IsOtpValid_Fail(t *testing.T) {
	var userOtpValidator validator.IUserOtpValidator
	userOtpValidator = validator.UserOtpValidator{}

	if userOtpValidator.IsOtpValid("12345", 6) {
		t.Fatal("expected false")
	}

	if userOtpValidator.IsOtpValid("abc3fa", 6) {
		t.Fatal("expected false")
	}
}



