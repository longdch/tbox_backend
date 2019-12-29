package helpers_test

import (
	"tbox_backend/internal/helpers"
	"testing"
)

func TestPhoneNumberRateLimiter_Limit(t *testing.T) {
	phoneNumberLimiters := helpers.NewPhoneNumberRateLimiters(1, 2)
	phoneNumber := "0967477177"
	limiter := phoneNumberLimiters.GetLimiter(phoneNumber)
	if !limiter.Allow() {
		t.Fatalf("expected true")
	}

	if !limiter.Allow() {
		t.Fatalf("expected true")
	}

	if limiter.Allow() {
		t.Fatalf("expected false")
	}
}

func TestPhoneNumberRateLimiter_DifferentPhoneNumber(t *testing.T) {
	phoneNumberLimiters := helpers.NewPhoneNumberRateLimiters(1, 1)
	if !phoneNumberLimiters.GetLimiter("0967467177").Allow() {
		t.Fatalf("expected true")
	}

	if phoneNumberLimiters.GetLimiter("0967467177").Allow() {
		t.Fatalf("expected false")
	}

	if !phoneNumberLimiters.GetLimiter("0988123123").Allow() {
		t.Fatalf("expected true")
	}
}
