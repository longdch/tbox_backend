package validator

import (
	"fmt"
	"regexp"
)

type IUserOtpValidator interface {
	IsOtpValid(otp string, size int) bool
}

type UserOtpValidator struct{}

func NewUserOtpValidator() *UserOtpValidator {
	return &UserOtpValidator{}
}

func (UserOtpValidator) IsOtpValid(otp string, size int) bool {
	regexPhoneNumber := fmt.Sprintf("^[0-9]{%d}$", size)
	regex := regexp.MustCompile(regexPhoneNumber)
	return regex.MatchString(otp)
}
