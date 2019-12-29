package validator

import (
	"regexp"
)

type IUserValidator interface {
	IsPhoneNumberValid(phoneNumber string) bool
}

type UserValidator struct{}

func NewUserValidator() *UserValidator {
	return &UserValidator{}
}

func (UserValidator) IsPhoneNumberValid(phoneNumber string) bool {
	regexPhoneNumber := `^(03|05|07|08|09)+([0-9]{8})$`
	regex := regexp.MustCompile(regexPhoneNumber)
	return regex.MatchString(phoneNumber)
}
