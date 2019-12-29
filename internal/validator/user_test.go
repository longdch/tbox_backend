package validator_test

import (
	"tbox_backend/internal/validator"
	"testing"
)

func TestUserValidator_IsPhoneNumberValid_Success(t *testing.T) {
	var userValidator validator.IUserValidator
	userValidator = validator.UserValidator{}
	phoneNumber := "0967499577"

	if !userValidator.IsPhoneNumberValid(phoneNumber) {
		t.Fatal("expected true")
	}
}

func TestUserValidator_IsPhoneNumberValid_Fail(t *testing.T) {
	var userValidator validator.IUserValidator
	userValidator = validator.UserValidator{}

	if userValidator.IsPhoneNumberValid("08674997777") {
		t.Fatal("expected false")
	}

	if userValidator.IsPhoneNumberValid("+84967471759") {
		t.Fatal("expected false")
	}

	if userValidator.IsPhoneNumberValid("0025185888") {
		t.Fatal("expected false")
	}
}
