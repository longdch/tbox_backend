package helpers_test

import (
	"tbox_backend/internal/helpers"
	"tbox_backend/internal/validator"
	"testing"
)

func TestUserOtpHelper_GenerateRandomOtp(t *testing.T) {
	userOtpHelper := helpers.NewUserOtpHelper()
	userOtpValidator := validator.NewUserOtpValidator()

	size := 6
	otp := userOtpHelper.GenerateRandomOtp(size)
	if !userOtpValidator.IsOtpValid(otp, size) {
		t.Fatalf("expected true")
	}
}
