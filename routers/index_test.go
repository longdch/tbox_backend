package routers_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"net/http"
	"net/http/httptest"
	"tbox_backend/internal/constants"
	"tbox_backend/internal/dto"
	"tbox_backend/internal/helpers"
	"tbox_backend/internal/validator"
	mockServices "tbox_backend/mock/services"
	"tbox_backend/routers"
	"testing"
)

func performRequest(r http.Handler, method, path string, body *bytes.Reader) *httptest.ResponseRecorder {
	req, _ := http.NewRequest(method, path, body)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	return w
}

func Test_GenerateOtp_Success(t *testing.T) {
	phoneNumber := "0967288123"
	gin.SetMode(gin.TestMode)
	router := gin.Default()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	userService := mockServices.NewMockIUserService(ctrl)
	userService.EXPECT().GenerateOtp(gomock.Eq(phoneNumber)).Return(nil)
	userValidator := validator.UserValidator{}
	phoneNumberLimiter := helpers.NewPhoneNumberRateLimiters(1, 1)
	r := routers.NewRouter(userService, userValidator, phoneNumberLimiter)

	r.IndexRouter(router)

	body := map[string]interface{}{
		"phone_number": phoneNumber,
	}

	postJson, _ := json.Marshal(body)
	w := performRequest(router, "POST", "/api/generate_otp", bytes.NewReader(postJson))

	if w.Code != http.StatusOK {
		t.Fatalf("expected status %d", http.StatusOK)
	}

	var response dto.GenerateOtpResponse
	err := json.Unmarshal([]byte(w.Body.String()), &response)
	if err != nil {
		t.Fatal(err)
	}

	if response.Message != "Success" {
		t.Fatalf("Expected success")
	}

	if response.Status != constants.SuccessStatus {
		t.Fatalf("Expected SuccessStatus")
	}
}

func Test_GenerateOtp_InvalidRequest(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.Default()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	userService := mockServices.NewMockIUserService(ctrl)
	userValidator := validator.UserValidator{}
	phoneNumberLimiter := helpers.NewPhoneNumberRateLimiters(1, 1)
	r := routers.NewRouter(userService, userValidator, phoneNumberLimiter)

	r.IndexRouter(router)
	w := performRequest(router, "POST", "/api/generate_otp", bytes.NewReader([]byte("random_text")))

	if w.Code != http.StatusOK {
		t.Fatalf("expected status %d", http.StatusOK)
	}

	var response dto.GenerateOtpResponse
	err := json.Unmarshal([]byte(w.Body.String()), &response)
	if err != nil {
		t.Fatal(err)
	}

	if response.Status != constants.InvalidRequestStatus {
		t.Fatalf("Expected InvalidRequestStatus")
	}
}

func Test_GenerateOtp_InvalidPhoneNumber(t *testing.T) {
	phoneNumber := "099967288123"
	router := gin.Default()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	userService := mockServices.NewMockIUserService(ctrl)
	userValidator := validator.UserValidator{}
	phoneNumberLimiter := helpers.NewPhoneNumberRateLimiters(1, 1)
	r := routers.NewRouter(userService, userValidator, phoneNumberLimiter)

	r.IndexRouter(router)

	body := map[string]interface{}{
		"phone_number": phoneNumber,
	}

	postJson, _ := json.Marshal(body)
	w := performRequest(router, "POST", "/api/generate_otp", bytes.NewReader(postJson))

	if w.Code != http.StatusOK {
		t.Fatalf("expected status %d", http.StatusOK)
	}

	var response dto.GenerateOtpResponse
	err := json.Unmarshal([]byte(w.Body.String()), &response)
	if err != nil {
		t.Fatal(err)
	}

	if response.Status != constants.InvalidRequestStatus {
		t.Fatalf("Expected InvalidRequestStatus")
	}
}

func Test_GenerateOtp_RateLimit(t *testing.T) {
	phoneNumber := "0967321412"
	gin.SetMode(gin.TestMode)
	router := gin.Default()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	userService := mockServices.NewMockIUserService(ctrl)
	userValidator := validator.UserValidator{}
	phoneNumberLimiter := helpers.NewPhoneNumberRateLimiters(1, 0)
	r := routers.NewRouter(userService, userValidator, phoneNumberLimiter)

	r.IndexRouter(router)

	body := map[string]interface{}{
		"phone_number": phoneNumber,
	}

	postJson, _ := json.Marshal(body)
	w := performRequest(router, "POST", "/api/generate_otp", bytes.NewReader(postJson))

	if w.Code != http.StatusOK {
		t.Fatalf("expected status %d", http.StatusOK)
	}

	var response dto.GenerateOtpResponse
	err := json.Unmarshal([]byte(w.Body.String()), &response)
	if err != nil {
		t.Fatal(err)
	}

	if response.Status != constants.TooManyRequestStatus {
		t.Fatalf("Expected TooManyRequestStatus")
	}
}

func Test_GenerateOtp_Error(t *testing.T) {
	phoneNumber := "0967288123"
	gin.SetMode(gin.TestMode)
	router := gin.Default()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	userService := mockServices.NewMockIUserService(ctrl)
	userService.EXPECT().GenerateOtp(gomock.Eq(phoneNumber)).Return(errors.New("Something went wrong "))
	userValidator := validator.UserValidator{}
	phoneNumberLimiter := helpers.NewPhoneNumberRateLimiters(1, 1)
	r := routers.NewRouter(userService, userValidator, phoneNumberLimiter)

	r.IndexRouter(router)

	body := map[string]interface{}{
		"phone_number": phoneNumber,
	}

	postJson, _ := json.Marshal(body)
	w := performRequest(router, "POST", "/api/generate_otp", bytes.NewReader(postJson))

	if w.Code != http.StatusOK {
		t.Fatalf("expected status %d", http.StatusOK)
	}

	var response dto.GenerateOtpResponse
	err := json.Unmarshal([]byte(w.Body.String()), &response)
	if err != nil {
		t.Fatal(err)
	}

	if response.Status != constants.SomethingWentWrongStatus {
		t.Fatalf("Expected SuccessStatus")
	}
}

func Test_ResendOtp_Success(t *testing.T) {
	phoneNumber := "0967288123"
	gin.SetMode(gin.TestMode)
	router := gin.Default()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	userService := mockServices.NewMockIUserService(ctrl)
	userService.EXPECT().ResendOtp(gomock.Eq(phoneNumber)).Return(nil)
	userValidator := validator.UserValidator{}
	phoneNumberLimiter := helpers.NewPhoneNumberRateLimiters(1, 1)
	r := routers.NewRouter(userService, userValidator, phoneNumberLimiter)

	r.IndexRouter(router)

	body := map[string]interface{}{
		"phone_number": phoneNumber,
	}

	postJson, _ := json.Marshal(body)
	w := performRequest(router, "POST", "/api/resend_otp", bytes.NewReader(postJson))

	if w.Code != http.StatusOK {
		t.Fatalf("expected status %d", http.StatusOK)
	}

	var response dto.GenerateOtpResponse
	err := json.Unmarshal([]byte(w.Body.String()), &response)
	if err != nil {
		t.Fatal(err)
	}

	if response.Message != "Success" {
		t.Fatalf("Expected success")
	}

	if response.Status != constants.SuccessStatus {
		t.Fatalf("Expected SuccessStatus")
	}
}

func Test_ResendOtp_InvalidRequest(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.Default()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	userService := mockServices.NewMockIUserService(ctrl)
	userValidator := validator.UserValidator{}
	phoneNumberLimiter := helpers.NewPhoneNumberRateLimiters(1, 1)
	r := routers.NewRouter(userService, userValidator, phoneNumberLimiter)

	r.IndexRouter(router)
	w := performRequest(router, "POST", "/api/resend_otp", bytes.NewReader([]byte("random_text")))

	if w.Code != http.StatusOK {
		t.Fatalf("expected status %d", http.StatusOK)
	}

	var response dto.GenerateOtpResponse
	err := json.Unmarshal([]byte(w.Body.String()), &response)
	if err != nil {
		t.Fatal(err)
	}

	if response.Status != constants.InvalidRequestStatus {
		t.Fatalf("Expected InvalidRequestStatus")
	}
}

func Test_ResendOtp_InvalidPhoneNumber(t *testing.T) {
	phoneNumber := "099967288123"
	gin.SetMode(gin.TestMode)
	router := gin.Default()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	userService := mockServices.NewMockIUserService(ctrl)
	userValidator := validator.UserValidator{}
	phoneNumberLimiter := helpers.NewPhoneNumberRateLimiters(1, 1)
	r := routers.NewRouter(userService, userValidator, phoneNumberLimiter)

	r.IndexRouter(router)

	body := map[string]interface{}{
		"phone_number": phoneNumber,
	}

	postJson, _ := json.Marshal(body)
	w := performRequest(router, "POST", "/api/resend_otp", bytes.NewReader(postJson))

	if w.Code != http.StatusOK {
		t.Fatalf("expected status %d", http.StatusOK)
	}

	var response dto.GenerateOtpResponse
	err := json.Unmarshal([]byte(w.Body.String()), &response)
	if err != nil {
		t.Fatal(err)
	}

	if response.Status != constants.InvalidRequestStatus {
		t.Fatalf("Expected InvalidRequestStatus")
	}
}

func Test_Resend_RateLimit(t *testing.T) {
	phoneNumber := "0967321412"
	gin.SetMode(gin.TestMode)
	router := gin.Default()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	userService := mockServices.NewMockIUserService(ctrl)
	userValidator := validator.UserValidator{}
	phoneNumberLimiter := helpers.NewPhoneNumberRateLimiters(1, 0)
	r := routers.NewRouter(userService, userValidator, phoneNumberLimiter)

	r.IndexRouter(router)

	body := map[string]interface{}{
		"phone_number": phoneNumber,
	}

	postJson, _ := json.Marshal(body)
	w := performRequest(router, "POST", "/api/resend_otp", bytes.NewReader(postJson))

	if w.Code != http.StatusOK {
		t.Fatalf("expected status %d", http.StatusOK)
	}

	var response dto.GenerateOtpResponse
	err := json.Unmarshal([]byte(w.Body.String()), &response)
	if err != nil {
		t.Fatal(err)
	}

	if response.Status != constants.TooManyRequestStatus {
		t.Fatalf("Expected TooManyRequestStatus")
	}
}

func Test_Resend_Error(t *testing.T) {
	phoneNumber := "0967288123"
	gin.SetMode(gin.TestMode)
	router := gin.Default()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	userService := mockServices.NewMockIUserService(ctrl)
	userService.EXPECT().ResendOtp(gomock.Eq(phoneNumber)).Return(errors.New("Something went wrong "))
	userValidator := validator.UserValidator{}
	phoneNumberLimiter := helpers.NewPhoneNumberRateLimiters(1, 1)
	r := routers.NewRouter(userService, userValidator, phoneNumberLimiter)

	r.IndexRouter(router)

	body := map[string]interface{}{
		"phone_number": phoneNumber,
	}

	postJson, _ := json.Marshal(body)
	w := performRequest(router, "POST", "/api/resend_otp", bytes.NewReader(postJson))

	if w.Code != http.StatusOK {
		t.Fatalf("expected status %d", http.StatusOK)
	}

	var response dto.GenerateOtpResponse
	err := json.Unmarshal([]byte(w.Body.String()), &response)
	if err != nil {
		t.Fatal(err)
	}

	if response.Status != constants.SomethingWentWrongStatus {
		t.Fatalf("Expected SuccessStatus")
	}
}

func Test_Login_Success(t *testing.T) {
	phoneNumber := "0967288123"
	gin.SetMode(gin.TestMode)
	router := gin.Default()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	userService := mockServices.NewMockIUserService(ctrl)
	userService.EXPECT().Login(gomock.Eq(phoneNumber), gomock.Any()).Return("tokentest", nil)
	userValidator := validator.UserValidator{}
	phoneNumberLimiter := helpers.NewPhoneNumberRateLimiters(0, 0)
	r := routers.NewRouter(userService, userValidator, phoneNumberLimiter)

	r.IndexRouter(router)
	body := map[string]interface{}{
		"phone_number": phoneNumber,
	}

	postJson, _ := json.Marshal(body)
	w := performRequest(router, "POST", "/api/login", bytes.NewReader(postJson))

	if w.Code != http.StatusOK {
		t.Fatalf("expected status %d", http.StatusOK)
	}

	var response dto.LoginResponse
	err := json.Unmarshal([]byte(w.Body.String()), &response)
	if err != nil {
		t.Fatal(err)
	}

	if response.Message != "Success" {
		t.Fatalf("Expected success")
	}

	if response.Status != constants.SuccessStatus {
		t.Fatalf("Expected SuccessStatus")
	}

	if response.Token != "tokentest" {
		t.Fatalf("Expected success")
	}
}

func Test_Login_InvalidRequest(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.Default()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	userService := mockServices.NewMockIUserService(ctrl)
	userValidator := validator.UserValidator{}
	phoneNumberLimiter := helpers.NewPhoneNumberRateLimiters(0, 0)
	r := routers.NewRouter(userService, userValidator, phoneNumberLimiter)

	r.IndexRouter(router)
	w := performRequest(router, "POST", "/api/login", bytes.NewReader([]byte("random_text")))

	if w.Code != http.StatusOK {
		t.Fatalf("expected status %d", http.StatusOK)
	}

	var response dto.LoginResponse
	err := json.Unmarshal([]byte(w.Body.String()), &response)
	if err != nil {
		t.Fatal(err)
	}

	if response.Status != constants.InvalidRequestStatus {
		t.Fatalf("Expected InvalidRequestStatus")
	}
}

func Test_Login_Error(t *testing.T) {
	phoneNumber := "0967288123"
	gin.SetMode(gin.TestMode)
	router := gin.Default()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	userService := mockServices.NewMockIUserService(ctrl)
	userService.EXPECT().Login(gomock.Eq(phoneNumber), gomock.Any()).Return("", errors.New("Something went wrong "))
	userValidator := validator.UserValidator{}
	phoneNumberLimiter := helpers.NewPhoneNumberRateLimiters(1, 1)
	r := routers.NewRouter(userService, userValidator, phoneNumberLimiter)

	r.IndexRouter(router)

	body := map[string]interface{}{
		"phone_number": phoneNumber,
	}

	postJson, _ := json.Marshal(body)
	w := performRequest(router, "POST", "/api/login", bytes.NewReader(postJson))

	if w.Code != http.StatusOK {
		t.Fatalf("expected status %d", http.StatusOK)
	}

	var response dto.LoginResponse
	err := json.Unmarshal([]byte(w.Body.String()), &response)
	if err != nil {
		t.Fatal(err)
	}

	if response.Status != constants.SomethingWentWrongStatus {
		t.Fatalf("Expected SuccessStatus")
	}
}
