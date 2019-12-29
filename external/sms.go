package external

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

type ISmsService interface {
	SendOtp(phoneNumber string, otp string) error
}

type SmsService struct {
	url string
}

func NewSmsService(url string) *SmsService {
	return &SmsService{url: url}
}

type SmsRequest struct {
	PhoneNumber string `json:"phone_number"`
	Content     string `json:"content"`
}

func (r SmsRequest) String() string {
	text, _ := json.Marshal(r)
	return string(text)
}

func (s SmsService) SendOtp(phoneNumber string, otp string) error {
	request := &SmsRequest{
		PhoneNumber: phoneNumber,
		Content:     fmt.Sprintf("Your OTP is: %s", otp),
	}

	log.Println("----- SMS message -----")
	log.Println(request)
	log.Println("-----------------------")

	buf := new(bytes.Buffer)
	_ = json.NewEncoder(buf).Encode(request)
	req, _ := http.NewRequest("POST", s.url, buf)
	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return err
	}

	defer func() {
		_ = res.Body.Close()
	}()

	return nil
}
