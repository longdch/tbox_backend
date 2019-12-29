package dto

type GenerateOtpRequest struct {
	PhoneNumber string `json:"phone_number"`
}

type LoginRequest struct {
	PhoneNumber string `json:"phone_number"`
	Otp         string `json:"otp"`
}
