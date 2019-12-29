package errors

import (
	"fmt"
)

type VerifiedPhoneNumberError struct {
	PhoneNumber string
}

func (e VerifiedPhoneNumberError) Error() string {
	return fmt.Sprintf("Phone number %s is already verified ", e.PhoneNumber)
}

type NotExistsPhoneNumberError struct {
	PhoneNumber string
}

func (e NotExistsPhoneNumberError) Error() string {
	return fmt.Sprintf("Phone number %s is not found ", e.PhoneNumber)
}

type InvalidPhoneNumberError struct {
	PhoneNumber string
}

func (e InvalidPhoneNumberError) Error() string {
	return fmt.Sprintf("Phone number %s is invalid ", e.PhoneNumber)
}

type GeneratedOtpError struct {
}

func (e GeneratedOtpError) Error() string {
	return "OTP has been generated "
}

type InvalidOtpError struct {
	Otp string
}

func (e InvalidOtpError) Error() string {
	return fmt.Sprintf("OTP %s is invalid ", e.Otp)
}

type IncorrectOtpError struct {
	Otp string
}

func (e IncorrectOtpError) Error() string {
	return fmt.Sprintf("OTP %s is incorrect ", e.Otp)
}

type ExpiredOtpError struct {
	Otp string
}

func (e ExpiredOtpError) Error() string {
	return fmt.Sprintf("OTP %s is expired ", e.Otp)
}
