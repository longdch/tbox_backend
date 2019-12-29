package dto_test

import (
	"tbox_backend/internal/dto"
	"testing"
)

func TestNewGenerateOtpResponse(t *testing.T) {
	generateOtpResponse := dto.NewGenerateOtpResponse(100, "test")
	if generateOtpResponse.Status != 100 {
		t.Fatalf("expected status: 100")
	}

	if generateOtpResponse.Message != "test" {
		t.Fatalf("expected message: test")
	}
}

func TestNewLoginResponse(t *testing.T) {
	loginResponse := dto.NewLoginResponse(100, "test", "abc")
	if loginResponse.Status != 100 {
		t.Fatalf("expected status: 100")
	}

	if loginResponse.Message != "test" {
		t.Fatalf("expected message: test")
	}

	if loginResponse.Token != "abc" {
		t.Fatalf("expected token: abc")
	}
}
