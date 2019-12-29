package services_test

import (
	"errors"
	"github.com/golang/mock/gomock"
	"tbox_backend/config"
	"tbox_backend/internal/constants"
	"tbox_backend/internal/dto"
	e "tbox_backend/internal/errors"
	"tbox_backend/internal/helpers"
	"tbox_backend/internal/services"
	"tbox_backend/internal/validator"
	mockExternal "tbox_backend/mock/external"
	mockStores "tbox_backend/mock/stores"
	"testing"
	"time"
)

func TestUserService_GenerateOtp_Success_FirstTime(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	phoneNumber := "0961234567"
	userStore := mockStores.NewMockIUserStore(ctrl)
	userStore.EXPECT().GetByPhoneNumber(gomock.Eq(phoneNumber)).Return(nil, false, nil)
	userID := 1
	userStore.EXPECT().Save(gomock.Any()).Do(func(user *dto.User) {
		user.ID = userID
	}).Return(nil)

	userOtpStore := mockStores.NewMockIUserOtpStore(ctrl)
	userOtpStore.EXPECT().GetByUserID(gomock.Eq(userID)).Return(dto.UserOtp{}, false, nil)
	userOtpStore.EXPECT().Save(gomock.Any()).Return(nil)

	smsService := mockExternal.NewMockISmsService(ctrl)
	smsService.EXPECT().SendOtp(gomock.Eq(phoneNumber), gomock.Any()).Return(nil)

	userValidator := validator.NewUserValidator()
	userOtpValidator := validator.NewUserOtpValidator()
	userOtpHelper := helpers.NewUserOtpHelper()
	userHelper := helpers.NewUserHelper("")

	cfg := config.Config{
		Base:                 config.Base{},
		MySQL:                config.MySQL{},
		PhoneNumberRateLimit: config.PhoneNumberRateLimit{},
		Otp:                  config.Otp{},
		SmsService:           config.SmsService{},
		Token:                config.Token{},
	}

	userService := services.NewUserService(
		cfg,
		smsService,
		userValidator,
		userOtpValidator,
		userOtpHelper,
		userHelper,
		userStore,
		userOtpStore,
	)

	err := userService.GenerateOtp(phoneNumber)
	if err != nil {
		t.Fatalf("expected nil")
	}
}

func TestUserService_GenerateOtp_Success_OtpExpired(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	phoneNumber := "0961234567"
	now := time.Now().UTC()
	tm := now.Add(-1 * time.Duration(61) * time.Second)

	userDto := &dto.User{
		ID:          1,
		PhoneNumber: phoneNumber,
		Status:      constants.UserInitStatus,
		CreatedAt:   tm,
		UpdatedAt:   tm,
	}

	userStore := mockStores.NewMockIUserStore(ctrl)
	userStore.EXPECT().GetByPhoneNumber(gomock.Eq(phoneNumber)).Return(userDto, true, nil)

	userOtpStore := mockStores.NewMockIUserOtpStore(ctrl)
	userOtpDto := dto.UserOtp{
		ID:        2,
		UserID:    1,
		Otp:       "123456",
		CreatedAt: tm,
		UpdatedAt: tm,
	}

	userOtpStore.EXPECT().GetByUserID(gomock.Eq(userDto.ID)).Return(userOtpDto, true, nil)
	userOtpStore.EXPECT().UpdateOtp(gomock.Any()).Return(nil)

	smsService := mockExternal.NewMockISmsService(ctrl)
	smsService.EXPECT().SendOtp(gomock.Eq(phoneNumber), gomock.Any()).Return(nil)

	userValidator := validator.NewUserValidator()
	userOtpValidator := validator.NewUserOtpValidator()
	userOtpHelper := helpers.NewUserOtpHelper()
	userHelper := helpers.NewUserHelper("")

	cfg := config.Config{
		Base:                 config.Base{},
		MySQL:                config.MySQL{},
		PhoneNumberRateLimit: config.PhoneNumberRateLimit{},
		Otp:                  config.Otp{},
		SmsService:           config.SmsService{},
		Token:                config.Token{},
	}

	cfg.Otp.ExpiredTime = 60
	userService := services.NewUserService(
		cfg,
		smsService,
		userValidator,
		userOtpValidator,
		userOtpHelper,
		userHelper,
		userStore,
		userOtpStore,
	)

	err := userService.GenerateOtp(phoneNumber)
	if err != nil {
		t.Fatalf("expected nil")
	}
}

func TestUserService_GenerateOtp_GetByPhoneNumber_Error(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	phoneNumber := "0961234567"
	userStore := mockStores.NewMockIUserStore(ctrl)
	expectedError := errors.New("Too many request ")
	userStore.EXPECT().GetByPhoneNumber(gomock.Eq(phoneNumber)).Return(nil, false, expectedError)

	userOtpStore := mockStores.NewMockIUserOtpStore(ctrl)
	smsService := mockExternal.NewMockISmsService(ctrl)
	userValidator := validator.NewUserValidator()
	userOtpValidator := validator.NewUserOtpValidator()
	userOtpHelper := helpers.NewUserOtpHelper()
	userHelper := helpers.NewUserHelper("")

	cfg := config.Config{
		Base:                 config.Base{},
		MySQL:                config.MySQL{},
		PhoneNumberRateLimit: config.PhoneNumberRateLimit{},
		Otp:                  config.Otp{},
		SmsService:           config.SmsService{},
		Token:                config.Token{},
	}

	userService := services.NewUserService(
		cfg,
		smsService,
		userValidator,
		userOtpValidator,
		userOtpHelper,
		userHelper,
		userStore,
		userOtpStore,
	)

	err := userService.GenerateOtp(phoneNumber)
	if err == nil || err.Error() != expectedError.Error() {
		t.Fatalf("expected err: %v", err)
	}
}

func TestUserService_GenerateOtp_SaveUser_Error(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	phoneNumber := "0961234567"
	userStore := mockStores.NewMockIUserStore(ctrl)
	expectedError := errors.New("Too many request ")
	userStore.EXPECT().GetByPhoneNumber(gomock.Eq(phoneNumber)).Return(nil, false, nil)
	userStore.EXPECT().Save(gomock.Any()).Return(expectedError)

	userOtpStore := mockStores.NewMockIUserOtpStore(ctrl)
	smsService := mockExternal.NewMockISmsService(ctrl)
	userValidator := validator.NewUserValidator()
	userOtpValidator := validator.NewUserOtpValidator()
	userOtpHelper := helpers.NewUserOtpHelper()
	userHelper := helpers.NewUserHelper("")

	cfg := config.Config{
		Base:                 config.Base{},
		MySQL:                config.MySQL{},
		PhoneNumberRateLimit: config.PhoneNumberRateLimit{},
		Otp:                  config.Otp{},
		SmsService:           config.SmsService{},
		Token:                config.Token{},
	}

	userService := services.NewUserService(
		cfg,
		smsService,
		userValidator,
		userOtpValidator,
		userOtpHelper,
		userHelper,
		userStore,
		userOtpStore,
	)

	err := userService.GenerateOtp(phoneNumber)
	if err == nil || err.Error() != expectedError.Error() {
		t.Fatalf("expected err: %v", expectedError)
	}
}

func TestUserService_GenerateOtp_Verified(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	phoneNumber := "0961234567"
	userStore := mockStores.NewMockIUserStore(ctrl)
	expectedError := e.VerifiedPhoneNumberError{PhoneNumber: phoneNumber}

	now := time.Now()
	userDto := &dto.User{
		ID:          1,
		PhoneNumber: phoneNumber,
		Status:      constants.UserVerifiedStatus,
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	userStore.EXPECT().GetByPhoneNumber(gomock.Eq(phoneNumber)).Return(userDto, true, nil)

	userOtpStore := mockStores.NewMockIUserOtpStore(ctrl)
	smsService := mockExternal.NewMockISmsService(ctrl)
	userValidator := validator.NewUserValidator()
	userOtpValidator := validator.NewUserOtpValidator()
	userOtpHelper := helpers.NewUserOtpHelper()
	userHelper := helpers.NewUserHelper("")

	cfg := config.Config{
		Base:                 config.Base{},
		MySQL:                config.MySQL{},
		PhoneNumberRateLimit: config.PhoneNumberRateLimit{},
		Otp:                  config.Otp{},
		SmsService:           config.SmsService{},
		Token:                config.Token{},
	}

	userService := services.NewUserService(
		cfg,
		smsService,
		userValidator,
		userOtpValidator,
		userOtpHelper,
		userHelper,
		userStore,
		userOtpStore,
	)

	err := userService.GenerateOtp(phoneNumber)
	if err == nil || err.Error() != expectedError.Error() {
		t.Fatalf("expected err: %v", expectedError)
	}
}

func TestUserService_GenerateOtp_GetUserOtp_Error(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	phoneNumber := "0961234567"
	now := time.Now().UTC()
	tm := now.Add(-1 * time.Duration(61) * time.Second)

	userDto := &dto.User{
		ID:          1,
		PhoneNumber: phoneNumber,
		Status:      constants.UserInitStatus,
		CreatedAt:   tm,
		UpdatedAt:   tm,
	}

	userStore := mockStores.NewMockIUserStore(ctrl)
	userStore.EXPECT().GetByPhoneNumber(gomock.Eq(phoneNumber)).Return(userDto, true, nil)

	userOtpStore := mockStores.NewMockIUserOtpStore(ctrl)
	expectedError := errors.New("Too many request ")
	userOtpStore.EXPECT().GetByUserID(gomock.Eq(userDto.ID)).Return(dto.UserOtp{}, false, expectedError)

	smsService := mockExternal.NewMockISmsService(ctrl)
	userValidator := validator.NewUserValidator()
	userOtpValidator := validator.NewUserOtpValidator()
	userOtpHelper := helpers.NewUserOtpHelper()
	userHelper := helpers.NewUserHelper("")

	cfg := config.Config{
		Base:                 config.Base{},
		MySQL:                config.MySQL{},
		PhoneNumberRateLimit: config.PhoneNumberRateLimit{},
		Otp:                  config.Otp{},
		SmsService:           config.SmsService{},
		Token:                config.Token{},
	}

	cfg.Otp.ExpiredTime = 60
	userService := services.NewUserService(
		cfg,
		smsService,
		userValidator,
		userOtpValidator,
		userOtpHelper,
		userHelper,
		userStore,
		userOtpStore,
	)

	err := userService.GenerateOtp(phoneNumber)
	if err == nil || err.Error() != expectedError.Error() {
		t.Fatalf("expected err: %v", expectedError)
	}
}

func TestUserService_GenerateOtp_Success_OtpNotExpired(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	phoneNumber := "0961234567"
	now := time.Now().UTC()
	tm := now.Add(-1 * time.Duration(30) * time.Second)

	userDto := &dto.User{
		ID:          1,
		PhoneNumber: phoneNumber,
		Status:      constants.UserInitStatus,
		CreatedAt:   tm,
		UpdatedAt:   tm,
	}

	userStore := mockStores.NewMockIUserStore(ctrl)
	userStore.EXPECT().GetByPhoneNumber(gomock.Eq(phoneNumber)).Return(userDto, true, nil)

	userOtpStore := mockStores.NewMockIUserOtpStore(ctrl)
	userOtpDto := dto.UserOtp{
		ID:        2,
		UserID:    1,
		Otp:       "123456",
		CreatedAt: tm,
		UpdatedAt: tm,
	}

	userOtpStore.EXPECT().GetByUserID(gomock.Eq(userDto.ID)).Return(userOtpDto, true, nil)

	smsService := mockExternal.NewMockISmsService(ctrl)
	userValidator := validator.NewUserValidator()
	userOtpValidator := validator.NewUserOtpValidator()
	userOtpHelper := helpers.NewUserOtpHelper()
	userHelper := helpers.NewUserHelper("")

	cfg := config.Config{
		Base:                 config.Base{},
		MySQL:                config.MySQL{},
		PhoneNumberRateLimit: config.PhoneNumberRateLimit{},
		Otp:                  config.Otp{},
		SmsService:           config.SmsService{},
		Token:                config.Token{},
	}

	cfg.Otp.ExpiredTime = 60
	userService := services.NewUserService(
		cfg,
		smsService,
		userValidator,
		userOtpValidator,
		userOtpHelper,
		userHelper,
		userStore,
		userOtpStore,
	)

	err := userService.GenerateOtp(phoneNumber)
	expectedError := e.GeneratedOtpError{}
	if err == nil || err.Error() != expectedError.Error() {
		t.Fatalf("expected error %v", expectedError)
	}
}

func TestUserService_GenerateOtp_UpdateOtp_Error(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	phoneNumber := "0961234567"
	now := time.Now().UTC()
	tm := now.Add(-1 * time.Duration(61) * time.Second)

	userDto := &dto.User{
		ID:          1,
		PhoneNumber: phoneNumber,
		Status:      constants.UserInitStatus,
		CreatedAt:   tm,
		UpdatedAt:   tm,
	}

	userStore := mockStores.NewMockIUserStore(ctrl)
	userStore.EXPECT().GetByPhoneNumber(gomock.Eq(phoneNumber)).Return(userDto, true, nil)

	userOtpStore := mockStores.NewMockIUserOtpStore(ctrl)
	userOtpDto := dto.UserOtp{
		ID:        2,
		UserID:    1,
		Otp:       "123456",
		CreatedAt: tm,
		UpdatedAt: tm,
	}

	userOtpStore.EXPECT().GetByUserID(gomock.Eq(userDto.ID)).Return(userOtpDto, true, nil)
	expectedError := errors.New("Too many request ")
	userOtpStore.EXPECT().UpdateOtp(gomock.Any()).Return(expectedError)

	smsService := mockExternal.NewMockISmsService(ctrl)
	userValidator := validator.NewUserValidator()
	userOtpValidator := validator.NewUserOtpValidator()
	userOtpHelper := helpers.NewUserOtpHelper()
	userHelper := helpers.NewUserHelper("")

	cfg := config.Config{
		Base:                 config.Base{},
		MySQL:                config.MySQL{},
		PhoneNumberRateLimit: config.PhoneNumberRateLimit{},
		Otp:                  config.Otp{},
		SmsService:           config.SmsService{},
		Token:                config.Token{},
	}

	cfg.Otp.ExpiredTime = 60
	userService := services.NewUserService(
		cfg,
		smsService,
		userValidator,
		userOtpValidator,
		userOtpHelper,
		userHelper,
		userStore,
		userOtpStore,
	)

	err := userService.GenerateOtp(phoneNumber)
	if err == nil || err.Error() != expectedError.Error() {
		t.Fatalf("expected error %v", expectedError)
	}
}

func TestUserService_GenerateOtp_SendSMS_Error(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	phoneNumber := "0961234567"
	now := time.Now().UTC()
	tm := now.Add(-1 * time.Duration(61) * time.Second)

	userDto := &dto.User{
		ID:          1,
		PhoneNumber: phoneNumber,
		Status:      constants.UserInitStatus,
		CreatedAt:   tm,
		UpdatedAt:   tm,
	}

	userStore := mockStores.NewMockIUserStore(ctrl)
	userStore.EXPECT().GetByPhoneNumber(gomock.Eq(phoneNumber)).Return(userDto, true, nil)

	userOtpStore := mockStores.NewMockIUserOtpStore(ctrl)
	userOtpDto := dto.UserOtp{
		ID:        2,
		UserID:    1,
		Otp:       "123456",
		CreatedAt: tm,
		UpdatedAt: tm,
	}

	userOtpStore.EXPECT().GetByUserID(gomock.Eq(userDto.ID)).Return(userOtpDto, true, nil)
	userOtpStore.EXPECT().UpdateOtp(gomock.Any()).Return(nil)

	smsService := mockExternal.NewMockISmsService(ctrl)
	smsService.EXPECT().SendOtp(gomock.Eq(phoneNumber), gomock.Any()).Return(errors.New("Nothing "))

	userValidator := validator.NewUserValidator()
	userOtpValidator := validator.NewUserOtpValidator()
	userOtpHelper := helpers.NewUserOtpHelper()
	userHelper := helpers.NewUserHelper("")

	cfg := config.Config{
		Base:                 config.Base{},
		MySQL:                config.MySQL{},
		PhoneNumberRateLimit: config.PhoneNumberRateLimit{},
		Otp:                  config.Otp{},
		SmsService:           config.SmsService{},
		Token:                config.Token{},
	}

	cfg.Otp.ExpiredTime = 60
	userService := services.NewUserService(
		cfg,
		smsService,
		userValidator,
		userOtpValidator,
		userOtpHelper,
		userHelper,
		userStore,
		userOtpStore,
	)

	err := userService.GenerateOtp(phoneNumber)
	if err != nil {
		t.Fatalf("expected nil")
	}
}

func TestUserService_GenerateOtp_SaveOtp_Error(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	phoneNumber := "0961234567"
	userStore := mockStores.NewMockIUserStore(ctrl)
	userStore.EXPECT().GetByPhoneNumber(gomock.Eq(phoneNumber)).Return(nil, false, nil)
	userID := 1
	userStore.EXPECT().Save(gomock.Any()).Do(func(user *dto.User) {
		user.ID = userID
	}).Return(nil)

	userOtpStore := mockStores.NewMockIUserOtpStore(ctrl)
	userOtpStore.EXPECT().GetByUserID(gomock.Eq(userID)).Return(dto.UserOtp{}, false, nil)
	expectedError := errors.New("Too many request ")
	userOtpStore.EXPECT().Save(gomock.Any()).Return(expectedError)

	smsService := mockExternal.NewMockISmsService(ctrl)
	userValidator := validator.NewUserValidator()
	userOtpValidator := validator.NewUserOtpValidator()
	userOtpHelper := helpers.NewUserOtpHelper()
	userHelper := helpers.NewUserHelper("")

	cfg := config.Config{
		Base:                 config.Base{},
		MySQL:                config.MySQL{},
		PhoneNumberRateLimit: config.PhoneNumberRateLimit{},
		Otp:                  config.Otp{},
		SmsService:           config.SmsService{},
		Token:                config.Token{},
	}

	userService := services.NewUserService(
		cfg,
		smsService,
		userValidator,
		userOtpValidator,
		userOtpHelper,
		userHelper,
		userStore,
		userOtpStore,
	)

	err := userService.GenerateOtp(phoneNumber)
	if err == nil || err.Error() != expectedError.Error() {
		t.Fatalf("expected error %v", expectedError)
	}
}

func TestUserService_GenerateOtp_FirstTime_SendSMS_Error(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	phoneNumber := "0961234567"
	userStore := mockStores.NewMockIUserStore(ctrl)
	userStore.EXPECT().GetByPhoneNumber(gomock.Eq(phoneNumber)).Return(nil, false, nil)
	userID := 1
	userStore.EXPECT().Save(gomock.Any()).Do(func(user *dto.User) {
		user.ID = userID
	}).Return(nil)

	userOtpStore := mockStores.NewMockIUserOtpStore(ctrl)
	userOtpStore.EXPECT().GetByUserID(gomock.Eq(userID)).Return(dto.UserOtp{}, false, nil)
	userOtpStore.EXPECT().Save(gomock.Any()).Return(nil)

	smsService := mockExternal.NewMockISmsService(ctrl)
	smsService.EXPECT().SendOtp(gomock.Eq(phoneNumber), gomock.Any()).Return(errors.New("Nothing "))

	userValidator := validator.NewUserValidator()
	userOtpValidator := validator.NewUserOtpValidator()
	userOtpHelper := helpers.NewUserOtpHelper()
	userHelper := helpers.NewUserHelper("")

	cfg := config.Config{
		Base:                 config.Base{},
		MySQL:                config.MySQL{},
		PhoneNumberRateLimit: config.PhoneNumberRateLimit{},
		Otp:                  config.Otp{},
		SmsService:           config.SmsService{},
		Token:                config.Token{},
	}

	userService := services.NewUserService(
		cfg,
		smsService,
		userValidator,
		userOtpValidator,
		userOtpHelper,
		userHelper,
		userStore,
		userOtpStore,
	)

	err := userService.GenerateOtp(phoneNumber)
	if err != nil {
		t.Fatalf("expected nil")
	}
}

func TestUserService_ResendOtp_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	phoneNumber := "0961234567"
	now := time.Now().UTC()
	tm := now.Add(-1 * time.Duration(32) * time.Second)

	userDto := &dto.User{
		ID:          1,
		PhoneNumber: phoneNumber,
		Status:      constants.UserInitStatus,
		CreatedAt:   tm,
		UpdatedAt:   tm,
	}

	userStore := mockStores.NewMockIUserStore(ctrl)
	userStore.EXPECT().GetByPhoneNumber(gomock.Eq(phoneNumber)).Return(userDto, true, nil)

	userOtpStore := mockStores.NewMockIUserOtpStore(ctrl)
	userOtpDto := dto.UserOtp{
		ID:        2,
		UserID:    1,
		Otp:       "123456",
		CreatedAt: tm,
		UpdatedAt: tm,
	}

	userOtpStore.EXPECT().GetByUserID(gomock.Eq(userDto.ID)).Return(userOtpDto, true, nil)
	userOtpStore.EXPECT().UpdateOtp(gomock.Any()).Return(nil)

	smsService := mockExternal.NewMockISmsService(ctrl)
	smsService.EXPECT().SendOtp(gomock.Eq(phoneNumber), gomock.Any()).Return(nil)

	userValidator := validator.NewUserValidator()
	userOtpValidator := validator.NewUserOtpValidator()
	userOtpHelper := helpers.NewUserOtpHelper()
	userHelper := helpers.NewUserHelper("")

	cfg := config.Config{
		Base:                 config.Base{},
		MySQL:                config.MySQL{},
		PhoneNumberRateLimit: config.PhoneNumberRateLimit{},
		Otp:                  config.Otp{},
		SmsService:           config.SmsService{},
		Token:                config.Token{},
	}

	cfg.Otp.ExpiredTime = 60
	cfg.Otp.ResendWaitingTime = 30
	userService := services.NewUserService(
		cfg,
		smsService,
		userValidator,
		userOtpValidator,
		userOtpHelper,
		userHelper,
		userStore,
		userOtpStore,
	)

	err := userService.ResendOtp(phoneNumber)
	if err != nil {
		t.Fatalf("expected nil")
	}
}

func TestUserService_ResendOtp_GetByPhoneNumber_Error(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	phoneNumber := "0961234567"
	userStore := mockStores.NewMockIUserStore(ctrl)
	expectedError := errors.New("Too many request ")
	userStore.EXPECT().GetByPhoneNumber(gomock.Eq(phoneNumber)).Return(nil, false, expectedError)

	userOtpStore := mockStores.NewMockIUserOtpStore(ctrl)
	smsService := mockExternal.NewMockISmsService(ctrl)
	userValidator := validator.NewUserValidator()
	userOtpValidator := validator.NewUserOtpValidator()
	userOtpHelper := helpers.NewUserOtpHelper()
	userHelper := helpers.NewUserHelper("")

	cfg := config.Config{
		Base:                 config.Base{},
		MySQL:                config.MySQL{},
		PhoneNumberRateLimit: config.PhoneNumberRateLimit{},
		Otp:                  config.Otp{},
		SmsService:           config.SmsService{},
		Token:                config.Token{},
	}

	userService := services.NewUserService(
		cfg,
		smsService,
		userValidator,
		userOtpValidator,
		userOtpHelper,
		userHelper,
		userStore,
		userOtpStore,
	)

	err := userService.ResendOtp(phoneNumber)
	if err == nil || err.Error() != expectedError.Error() {
		t.Fatalf("expected error %v", expectedError)
	}
}

func TestUserService_ResendOtp_PhoneNumberNotExists(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	phoneNumber := "0961234567"
	userStore := mockStores.NewMockIUserStore(ctrl)
	userStore.EXPECT().GetByPhoneNumber(gomock.Eq(phoneNumber)).Return(nil, false, nil)

	userOtpStore := mockStores.NewMockIUserOtpStore(ctrl)
	smsService := mockExternal.NewMockISmsService(ctrl)
	userValidator := validator.NewUserValidator()
	userOtpValidator := validator.NewUserOtpValidator()
	userOtpHelper := helpers.NewUserOtpHelper()
	userHelper := helpers.NewUserHelper("")

	cfg := config.Config{
		Base:                 config.Base{},
		MySQL:                config.MySQL{},
		PhoneNumberRateLimit: config.PhoneNumberRateLimit{},
		Otp:                  config.Otp{},
		SmsService:           config.SmsService{},
		Token:                config.Token{},
	}

	userService := services.NewUserService(
		cfg,
		smsService,
		userValidator,
		userOtpValidator,
		userOtpHelper,
		userHelper,
		userStore,
		userOtpStore,
	)

	err := userService.ResendOtp(phoneNumber)
	expectedError := e.NotExistsPhoneNumberError{PhoneNumber: phoneNumber}
	if err == nil || err.Error() != expectedError.Error() {
		t.Fatalf("expected error %v", expectedError)
	}
}

func TestUserService_ResendOtp_PhoneNumber_Verified(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	phoneNumber := "0961234567"
	now := time.Now().UTC()
	tm := now.Add(-1 * time.Duration(32) * time.Second)

	userDto := &dto.User{
		ID:          1,
		PhoneNumber: phoneNumber,
		Status:      constants.UserVerifiedStatus,
		CreatedAt:   tm,
		UpdatedAt:   tm,
	}

	userStore := mockStores.NewMockIUserStore(ctrl)
	userStore.EXPECT().GetByPhoneNumber(gomock.Eq(phoneNumber)).Return(userDto, true, nil)

	userOtpStore := mockStores.NewMockIUserOtpStore(ctrl)
	smsService := mockExternal.NewMockISmsService(ctrl)
	userValidator := validator.NewUserValidator()
	userOtpValidator := validator.NewUserOtpValidator()
	userOtpHelper := helpers.NewUserOtpHelper()
	userHelper := helpers.NewUserHelper("")

	cfg := config.Config{
		Base:                 config.Base{},
		MySQL:                config.MySQL{},
		PhoneNumberRateLimit: config.PhoneNumberRateLimit{},
		Otp:                  config.Otp{},
		SmsService:           config.SmsService{},
		Token:                config.Token{},
	}

	userService := services.NewUserService(
		cfg,
		smsService,
		userValidator,
		userOtpValidator,
		userOtpHelper,
		userHelper,
		userStore,
		userOtpStore,
	)

	err := userService.ResendOtp(phoneNumber)
	expectedError := e.VerifiedPhoneNumberError{PhoneNumber: phoneNumber}
	if err == nil || err.Error() != expectedError.Error() {
		t.Fatalf("expected error %v", expectedError)
	}
}

func TestUserService_ResendOtp_GetByUserID_Error(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	phoneNumber := "0961234567"
	now := time.Now().UTC()
	userDto := &dto.User{
		ID:          1,
		PhoneNumber: phoneNumber,
		Status:      constants.UserInitStatus,
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	userStore := mockStores.NewMockIUserStore(ctrl)
	userStore.EXPECT().GetByPhoneNumber(gomock.Eq(phoneNumber)).Return(userDto, true, nil)

	userOtpStore := mockStores.NewMockIUserOtpStore(ctrl)
	expectedError := errors.New("Too many request ")
	userOtpStore.EXPECT().GetByUserID(gomock.Eq(userDto.ID)).Return(dto.UserOtp{}, false, expectedError)

	smsService := mockExternal.NewMockISmsService(ctrl)

	userValidator := validator.NewUserValidator()
	userOtpValidator := validator.NewUserOtpValidator()
	userOtpHelper := helpers.NewUserOtpHelper()
	userHelper := helpers.NewUserHelper("")

	cfg := config.Config{
		Base:                 config.Base{},
		MySQL:                config.MySQL{},
		PhoneNumberRateLimit: config.PhoneNumberRateLimit{},
		Otp:                  config.Otp{},
		SmsService:           config.SmsService{},
		Token:                config.Token{},
	}

	userService := services.NewUserService(
		cfg,
		smsService,
		userValidator,
		userOtpValidator,
		userOtpHelper,
		userHelper,
		userStore,
		userOtpStore,
	)

	err := userService.ResendOtp(phoneNumber)
	if err == nil || err.Error() != expectedError.Error() {
		t.Fatalf("expected error %v", expectedError)
	}
}

func TestUserService_ResendOtp_OtpNotExists(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	phoneNumber := "0961234567"
	now := time.Now().UTC()
	userDto := &dto.User{
		ID:          1,
		PhoneNumber: phoneNumber,
		Status:      constants.UserInitStatus,
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	userStore := mockStores.NewMockIUserStore(ctrl)
	userStore.EXPECT().GetByPhoneNumber(gomock.Eq(phoneNumber)).Return(userDto, true, nil)

	userOtpStore := mockStores.NewMockIUserOtpStore(ctrl)
	expectedError := errors.New("Could not resend OTP ")
	userOtpStore.EXPECT().GetByUserID(gomock.Eq(userDto.ID)).Return(dto.UserOtp{}, false, nil)

	smsService := mockExternal.NewMockISmsService(ctrl)
	userValidator := validator.NewUserValidator()
	userOtpValidator := validator.NewUserOtpValidator()
	userOtpHelper := helpers.NewUserOtpHelper()
	userHelper := helpers.NewUserHelper("")

	cfg := config.Config{
		Base:                 config.Base{},
		MySQL:                config.MySQL{},
		PhoneNumberRateLimit: config.PhoneNumberRateLimit{},
		Otp:                  config.Otp{},
		SmsService:           config.SmsService{},
		Token:                config.Token{},
	}

	userService := services.NewUserService(
		cfg,
		smsService,
		userValidator,
		userOtpValidator,
		userOtpHelper,
		userHelper,
		userStore,
		userOtpStore,
	)

	err := userService.ResendOtp(phoneNumber)
	if err == nil || err.Error() != expectedError.Error() {
		t.Fatalf("expected error %v", expectedError)
	}
}

func TestUserService_ResendOtp_OtpGenerated(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	phoneNumber := "0961234567"
	now := time.Now().UTC()
	tm := now.Add(-1 * time.Duration(10) * time.Second)

	userDto := &dto.User{
		ID:          1,
		PhoneNumber: phoneNumber,
		Status:      constants.UserInitStatus,
		CreatedAt:   tm,
		UpdatedAt:   tm,
	}

	userStore := mockStores.NewMockIUserStore(ctrl)
	userStore.EXPECT().GetByPhoneNumber(gomock.Eq(phoneNumber)).Return(userDto, true, nil)

	userOtpStore := mockStores.NewMockIUserOtpStore(ctrl)
	userOtpDto := dto.UserOtp{
		ID:        2,
		UserID:    1,
		Otp:       "123456",
		CreatedAt: tm,
		UpdatedAt: tm,
	}

	userOtpStore.EXPECT().GetByUserID(gomock.Eq(userDto.ID)).Return(userOtpDto, true, nil)

	smsService := mockExternal.NewMockISmsService(ctrl)
	userValidator := validator.NewUserValidator()
	userOtpValidator := validator.NewUserOtpValidator()
	userOtpHelper := helpers.NewUserOtpHelper()
	userHelper := helpers.NewUserHelper("")

	cfg := config.Config{
		Base:                 config.Base{},
		MySQL:                config.MySQL{},
		PhoneNumberRateLimit: config.PhoneNumberRateLimit{},
		Otp:                  config.Otp{},
		SmsService:           config.SmsService{},
		Token:                config.Token{},
	}

	cfg.Otp.ExpiredTime = 60
	cfg.Otp.ResendWaitingTime = 30
	userService := services.NewUserService(
		cfg,
		smsService,
		userValidator,
		userOtpValidator,
		userOtpHelper,
		userHelper,
		userStore,
		userOtpStore,
	)

	expectedError := e.GeneratedOtpError{}
	err := userService.ResendOtp(phoneNumber)
	if err == nil || err.Error() != expectedError.Error() {
		t.Fatalf("expected error %v", expectedError)
	}
}

func TestUserService_ResendOtp_UpdateOtp_Error(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	phoneNumber := "0961234567"
	now := time.Now().UTC()
	tm := now.Add(-1 * time.Duration(32) * time.Second)

	userDto := &dto.User{
		ID:          1,
		PhoneNumber: phoneNumber,
		Status:      constants.UserInitStatus,
		CreatedAt:   tm,
		UpdatedAt:   tm,
	}

	userStore := mockStores.NewMockIUserStore(ctrl)
	userStore.EXPECT().GetByPhoneNumber(gomock.Eq(phoneNumber)).Return(userDto, true, nil)

	userOtpStore := mockStores.NewMockIUserOtpStore(ctrl)
	userOtpDto := dto.UserOtp{
		ID:        2,
		UserID:    1,
		Otp:       "123456",
		CreatedAt: tm,
		UpdatedAt: tm,
	}

	userOtpStore.EXPECT().GetByUserID(gomock.Eq(userDto.ID)).Return(userOtpDto, true, nil)
	expectedError := errors.New("Too many request ")
	userOtpStore.EXPECT().UpdateOtp(gomock.Any()).Return(expectedError)

	smsService := mockExternal.NewMockISmsService(ctrl)
	userValidator := validator.NewUserValidator()
	userOtpValidator := validator.NewUserOtpValidator()
	userOtpHelper := helpers.NewUserOtpHelper()
	userHelper := helpers.NewUserHelper("")

	cfg := config.Config{
		Base:                 config.Base{},
		MySQL:                config.MySQL{},
		PhoneNumberRateLimit: config.PhoneNumberRateLimit{},
		Otp:                  config.Otp{},
		SmsService:           config.SmsService{},
		Token:                config.Token{},
	}

	cfg.Otp.ExpiredTime = 60
	cfg.Otp.ResendWaitingTime = 30
	userService := services.NewUserService(
		cfg,
		smsService,
		userValidator,
		userOtpValidator,
		userOtpHelper,
		userHelper,
		userStore,
		userOtpStore,
	)

	err := userService.ResendOtp(phoneNumber)
	if err == nil || err.Error() != expectedError.Error() {
		t.Fatalf("expected error %v", expectedError)
	}
}

func TestUserService_ResendOtp_SendSMS_Error(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	phoneNumber := "0961234567"
	now := time.Now().UTC()
	tm := now.Add(-1 * time.Duration(32) * time.Second)

	userDto := &dto.User{
		ID:          1,
		PhoneNumber: phoneNumber,
		Status:      constants.UserInitStatus,
		CreatedAt:   tm,
		UpdatedAt:   tm,
	}

	userStore := mockStores.NewMockIUserStore(ctrl)
	userStore.EXPECT().GetByPhoneNumber(gomock.Eq(phoneNumber)).Return(userDto, true, nil)

	userOtpStore := mockStores.NewMockIUserOtpStore(ctrl)
	userOtpDto := dto.UserOtp{
		ID:        2,
		UserID:    1,
		Otp:       "123456",
		CreatedAt: tm,
		UpdatedAt: tm,
	}

	userOtpStore.EXPECT().GetByUserID(gomock.Eq(userDto.ID)).Return(userOtpDto, true, nil)
	userOtpStore.EXPECT().UpdateOtp(gomock.Any()).Return(nil)

	smsService := mockExternal.NewMockISmsService(ctrl)
	smsService.EXPECT().SendOtp(gomock.Eq(phoneNumber), gomock.Any()).Return(errors.New("Nothing "))

	userValidator := validator.NewUserValidator()
	userOtpValidator := validator.NewUserOtpValidator()
	userOtpHelper := helpers.NewUserOtpHelper()
	userHelper := helpers.NewUserHelper("")

	cfg := config.Config{
		Base:                 config.Base{},
		MySQL:                config.MySQL{},
		PhoneNumberRateLimit: config.PhoneNumberRateLimit{},
		Otp:                  config.Otp{},
		SmsService:           config.SmsService{},
		Token:                config.Token{},
	}

	cfg.Otp.ExpiredTime = 60
	cfg.Otp.ResendWaitingTime = 30
	userService := services.NewUserService(
		cfg,
		smsService,
		userValidator,
		userOtpValidator,
		userOtpHelper,
		userHelper,
		userStore,
		userOtpStore,
	)

	err := userService.ResendOtp(phoneNumber)
	if err != nil {
		t.Fatalf("expected nil")
	}
}

func TestUserService_Login_InvalidPhoneNumber(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	phoneNumber := "10961234567"
	otp := "123456"

	userStore := mockStores.NewMockIUserStore(ctrl)
	userOtpStore := mockStores.NewMockIUserOtpStore(ctrl)
	smsService := mockExternal.NewMockISmsService(ctrl)
	userValidator := validator.NewUserValidator()
	userOtpValidator := validator.NewUserOtpValidator()
	userOtpHelper := helpers.NewUserOtpHelper()
	userHelper := helpers.NewUserHelper("abc")

	cfg := config.Config{
		Base:                 config.Base{},
		MySQL:                config.MySQL{},
		PhoneNumberRateLimit: config.PhoneNumberRateLimit{},
		Otp:                  config.Otp{},
		SmsService:           config.SmsService{},
		Token:                config.Token{},
	}

	cfg.Otp.ExpiredTime = 60
	userService := services.NewUserService(
		cfg,
		smsService,
		userValidator,
		userOtpValidator,
		userOtpHelper,
		userHelper,
		userStore,
		userOtpStore,
	)

	_, err := userService.Login(phoneNumber, otp)
	expectedError := e.InvalidPhoneNumberError{PhoneNumber: phoneNumber}
	if err == nil || err.Error() != expectedError.Error() {
		t.Fatalf("expect error %v", expectedError)
	}
}

func TestUserService_Login_GetPhoneNumberError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	phoneNumber := "0961234567"
	otp := "123456"

	userStore := mockStores.NewMockIUserStore(ctrl)
	expectedError := errors.New("Too many request ")
	userStore.EXPECT().GetByPhoneNumber(gomock.Eq(phoneNumber)).Return(nil, true, expectedError)

	userOtpStore := mockStores.NewMockIUserOtpStore(ctrl)
	smsService := mockExternal.NewMockISmsService(ctrl)
	userValidator := validator.NewUserValidator()
	userOtpValidator := validator.NewUserOtpValidator()
	userOtpHelper := helpers.NewUserOtpHelper()
	userHelper := helpers.NewUserHelper("abc")

	cfg := config.Config{
		Base:                 config.Base{},
		MySQL:                config.MySQL{},
		PhoneNumberRateLimit: config.PhoneNumberRateLimit{},
		Otp:                  config.Otp{},
		SmsService:           config.SmsService{},
		Token:                config.Token{},
	}

	cfg.Otp.ExpiredTime = 60
	userService := services.NewUserService(
		cfg,
		smsService,
		userValidator,
		userOtpValidator,
		userOtpHelper,
		userHelper,
		userStore,
		userOtpStore,
	)

	_, err := userService.Login(phoneNumber, otp)
	if err == nil || err.Error() != expectedError.Error() {
		t.Fatalf("expected error %v", expectedError)
	}
}

func TestUserService_Login_PhoneNumberNotExists(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	phoneNumber := "0961234567"
	otp := "123456"

	userStore := mockStores.NewMockIUserStore(ctrl)
	userStore.EXPECT().GetByPhoneNumber(gomock.Eq(phoneNumber)).Return(nil, false, nil)

	userOtpStore := mockStores.NewMockIUserOtpStore(ctrl)
	smsService := mockExternal.NewMockISmsService(ctrl)
	userValidator := validator.NewUserValidator()
	userOtpValidator := validator.NewUserOtpValidator()
	userOtpHelper := helpers.NewUserOtpHelper()
	userHelper := helpers.NewUserHelper("abc")

	cfg := config.Config{
		Base:                 config.Base{},
		MySQL:                config.MySQL{},
		PhoneNumberRateLimit: config.PhoneNumberRateLimit{},
		Otp:                  config.Otp{},
		SmsService:           config.SmsService{},
		Token:                config.Token{},
	}

	cfg.Otp.ExpiredTime = 60
	userService := services.NewUserService(
		cfg,
		smsService,
		userValidator,
		userOtpValidator,
		userOtpHelper,
		userHelper,
		userStore,
		userOtpStore,
	)

	_, err := userService.Login(phoneNumber, otp)
	expectedError := e.NotExistsPhoneNumberError{PhoneNumber: phoneNumber}
	if err == nil || err.Error() != expectedError.Error() {
		t.Fatalf("expect error %v", expectedError)
	}
}

func TestUserService_Login_PhoneNumberVerified(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	phoneNumber := "0961234567"
	otp := "123456"
	now := time.Now().UTC()
	tm := now.Add(-1 * time.Duration(32) * time.Second)

	userDto := &dto.User{
		ID:          1,
		PhoneNumber: phoneNumber,
		Status:      constants.UserVerifiedStatus,
		CreatedAt:   tm,
		UpdatedAt:   tm,
	}

	userStore := mockStores.NewMockIUserStore(ctrl)
	userStore.EXPECT().GetByPhoneNumber(gomock.Eq(phoneNumber)).Return(userDto, true, nil)

	userOtpStore := mockStores.NewMockIUserOtpStore(ctrl)
	smsService := mockExternal.NewMockISmsService(ctrl)
	userValidator := validator.NewUserValidator()
	userOtpValidator := validator.NewUserOtpValidator()
	userOtpHelper := helpers.NewUserOtpHelper()
	userHelper := helpers.NewUserHelper("abc")

	cfg := config.Config{
		Base:                 config.Base{},
		MySQL:                config.MySQL{},
		PhoneNumberRateLimit: config.PhoneNumberRateLimit{},
		Otp:                  config.Otp{},
		SmsService:           config.SmsService{},
		Token:                config.Token{},
	}

	cfg.Otp.ExpiredTime = 60
	userService := services.NewUserService(
		cfg,
		smsService,
		userValidator,
		userOtpValidator,
		userOtpHelper,
		userHelper,
		userStore,
		userOtpStore,
	)

	token, err := userService.Login(phoneNumber, otp)
	if err != nil {
		t.Fatalf("expected nil")
	}

	if token != "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjoxfQ.6dJuEd4gOyQ3aKNHDABXEwQrTYQgNPVDLd0ZcyIUw-Q" {
		t.Fatalf("wrong token")
	}
}

func TestUserService_Login_OtpInvalid(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	phoneNumber := "0961234567"
	otp := "1234"
	now := time.Now().UTC()
	tm := now.Add(-1 * time.Duration(32) * time.Second)

	userDto := &dto.User{
		ID:          1,
		PhoneNumber: phoneNumber,
		Status:      constants.UserInitStatus,
		CreatedAt:   tm,
		UpdatedAt:   tm,
	}

	userStore := mockStores.NewMockIUserStore(ctrl)
	userStore.EXPECT().GetByPhoneNumber(gomock.Eq(phoneNumber)).Return(userDto, true, nil)

	userOtpStore := mockStores.NewMockIUserOtpStore(ctrl)
	smsService := mockExternal.NewMockISmsService(ctrl)
	userValidator := validator.NewUserValidator()
	userOtpValidator := validator.NewUserOtpValidator()
	userOtpHelper := helpers.NewUserOtpHelper()
	userHelper := helpers.NewUserHelper("abc")

	cfg := config.Config{
		Base:                 config.Base{},
		MySQL:                config.MySQL{},
		PhoneNumberRateLimit: config.PhoneNumberRateLimit{},
		Otp:                  config.Otp{},
		SmsService:           config.SmsService{},
		Token:                config.Token{},
	}

	cfg.Otp.ExpiredTime = 60
	cfg.Otp.Size = 6
	userService := services.NewUserService(
		cfg,
		smsService,
		userValidator,
		userOtpValidator,
		userOtpHelper,
		userHelper,
		userStore,
		userOtpStore,
	)

	_, err := userService.Login(phoneNumber, otp)
	expectedError := e.InvalidOtpError{Otp: otp}
	if err == nil || err.Error() != expectedError.Error() {
		t.Fatalf("expect error %v", expectedError)
	}
}

func TestUserService_Login_GetByUserID_Error(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	phoneNumber := "0961234567"
	otp := "123456"
	now := time.Now().UTC()
	tm := now.Add(-1 * time.Duration(32) * time.Second)

	userDto := &dto.User{
		ID:          1,
		PhoneNumber: phoneNumber,
		Status:      constants.UserInitStatus,
		CreatedAt:   tm,
		UpdatedAt:   tm,
	}

	userStore := mockStores.NewMockIUserStore(ctrl)
	userStore.EXPECT().GetByPhoneNumber(gomock.Eq(phoneNumber)).Return(userDto, true, nil)

	userOtpStore := mockStores.NewMockIUserOtpStore(ctrl)
	expectedError := errors.New("Too many request ")
	userOtpStore.EXPECT().GetByUserID(gomock.Eq(userDto.ID)).Return(dto.UserOtp{}, true, expectedError)

	smsService := mockExternal.NewMockISmsService(ctrl)
	userValidator := validator.NewUserValidator()
	userOtpValidator := validator.NewUserOtpValidator()
	userOtpHelper := helpers.NewUserOtpHelper()
	userHelper := helpers.NewUserHelper("abc")

	cfg := config.Config{
		Base:                 config.Base{},
		MySQL:                config.MySQL{},
		PhoneNumberRateLimit: config.PhoneNumberRateLimit{},
		Otp:                  config.Otp{},
		SmsService:           config.SmsService{},
		Token:                config.Token{},
	}

	cfg.Otp.ExpiredTime = 60
	cfg.Otp.Size = 6
	userService := services.NewUserService(
		cfg,
		smsService,
		userValidator,
		userOtpValidator,
		userOtpHelper,
		userHelper,
		userStore,
		userOtpStore,
	)

	_, err := userService.Login(phoneNumber, otp)
	if err == nil || err.Error() != expectedError.Error() {
		t.Fatalf("expect error %v", expectedError)
	}
}

func TestUserService_Login_GetByUserID_NotExists(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	phoneNumber := "0961234567"
	otp := "123456"
	now := time.Now().UTC()
	tm := now.Add(-1 * time.Duration(32) * time.Second)

	userDto := &dto.User{
		ID:          1,
		PhoneNumber: phoneNumber,
		Status:      constants.UserInitStatus,
		CreatedAt:   tm,
		UpdatedAt:   tm,
	}

	userStore := mockStores.NewMockIUserStore(ctrl)
	userStore.EXPECT().GetByPhoneNumber(gomock.Eq(phoneNumber)).Return(userDto, true, nil)

	userOtpStore := mockStores.NewMockIUserOtpStore(ctrl)
	userOtpStore.EXPECT().GetByUserID(gomock.Eq(userDto.ID)).Return(dto.UserOtp{}, false, nil)

	smsService := mockExternal.NewMockISmsService(ctrl)
	userValidator := validator.NewUserValidator()
	userOtpValidator := validator.NewUserOtpValidator()
	userOtpHelper := helpers.NewUserOtpHelper()
	userHelper := helpers.NewUserHelper("abc")

	cfg := config.Config{
		Base:                 config.Base{},
		MySQL:                config.MySQL{},
		PhoneNumberRateLimit: config.PhoneNumberRateLimit{},
		Otp:                  config.Otp{},
		SmsService:           config.SmsService{},
		Token:                config.Token{},
	}

	cfg.Otp.ExpiredTime = 60
	cfg.Otp.Size = 6
	userService := services.NewUserService(
		cfg,
		smsService,
		userValidator,
		userOtpValidator,
		userOtpHelper,
		userHelper,
		userStore,
		userOtpStore,
	)

	_, err := userService.Login(phoneNumber, otp)
	expectedError := e.IncorrectOtpError{Otp: otp}
	if err == nil || err.Error() != expectedError.Error() {
		t.Fatalf("expect error %v", expectedError)
	}
}

func TestUserService_Login_Otp_Incorrect(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	phoneNumber := "0961234567"
	otp := "123457"
	now := time.Now().UTC()
	tm := now.Add(-1 * time.Duration(32) * time.Second)

	userDto := &dto.User{
		ID:          1,
		PhoneNumber: phoneNumber,
		Status:      constants.UserInitStatus,
		CreatedAt:   tm,
		UpdatedAt:   tm,
	}

	userStore := mockStores.NewMockIUserStore(ctrl)
	userStore.EXPECT().GetByPhoneNumber(gomock.Eq(phoneNumber)).Return(userDto, true, nil)

	userOtpStore := mockStores.NewMockIUserOtpStore(ctrl)
	userOtpDto := dto.UserOtp{
		ID:        2,
		UserID:    1,
		Otp:       "123456",
		CreatedAt: tm,
		UpdatedAt: tm,
	}

	userOtpStore.EXPECT().GetByUserID(gomock.Eq(userDto.ID)).Return(userOtpDto, true, nil)

	smsService := mockExternal.NewMockISmsService(ctrl)
	userValidator := validator.NewUserValidator()
	userOtpValidator := validator.NewUserOtpValidator()
	userOtpHelper := helpers.NewUserOtpHelper()
	userHelper := helpers.NewUserHelper("abc")

	cfg := config.Config{
		Base:                 config.Base{},
		MySQL:                config.MySQL{},
		PhoneNumberRateLimit: config.PhoneNumberRateLimit{},
		Otp:                  config.Otp{},
		SmsService:           config.SmsService{},
		Token:                config.Token{},
	}

	cfg.Otp.ExpiredTime = 60
	cfg.Otp.Size = 6
	userService := services.NewUserService(
		cfg,
		smsService,
		userValidator,
		userOtpValidator,
		userOtpHelper,
		userHelper,
		userStore,
		userOtpStore,
	)

	_, err := userService.Login(phoneNumber, otp)
	expectedError := e.IncorrectOtpError{Otp: otp}
	if err == nil || err.Error() != expectedError.Error() {
		t.Fatalf("expect error %v", expectedError)
	}
}

func TestUserService_Login_Otp_Expired(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	phoneNumber := "0961234567"
	otp := "123456"
	now := time.Now().UTC()
	tm := now.Add(-1 * time.Duration(70) * time.Second)

	userDto := &dto.User{
		ID:          1,
		PhoneNumber: phoneNumber,
		Status:      constants.UserInitStatus,
		CreatedAt:   tm,
		UpdatedAt:   tm,
	}

	userStore := mockStores.NewMockIUserStore(ctrl)
	userStore.EXPECT().GetByPhoneNumber(gomock.Eq(phoneNumber)).Return(userDto, true, nil)

	userOtpStore := mockStores.NewMockIUserOtpStore(ctrl)
	userOtpDto := dto.UserOtp{
		ID:        2,
		UserID:    1,
		Otp:       "123456",
		CreatedAt: tm,
		UpdatedAt: tm,
	}

	userOtpStore.EXPECT().GetByUserID(gomock.Eq(userDto.ID)).Return(userOtpDto, true, nil)

	smsService := mockExternal.NewMockISmsService(ctrl)
	userValidator := validator.NewUserValidator()
	userOtpValidator := validator.NewUserOtpValidator()
	userOtpHelper := helpers.NewUserOtpHelper()
	userHelper := helpers.NewUserHelper("abc")

	cfg := config.Config{
		Base:                 config.Base{},
		MySQL:                config.MySQL{},
		PhoneNumberRateLimit: config.PhoneNumberRateLimit{},
		Otp:                  config.Otp{},
		SmsService:           config.SmsService{},
		Token:                config.Token{},
	}

	cfg.Otp.ExpiredTime = 60
	cfg.Otp.Size = 6
	userService := services.NewUserService(
		cfg,
		smsService,
		userValidator,
		userOtpValidator,
		userOtpHelper,
		userHelper,
		userStore,
		userOtpStore,
	)

	_, err := userService.Login(phoneNumber, otp)
	expectedError := e.ExpiredOtpError{Otp: otp}
	if err == nil || err.Error() != expectedError.Error() {
		t.Fatalf("expect error %v", expectedError)
	}
}

func TestUserService_Login_UpdateUser_Error(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	phoneNumber := "0961234567"
	otp := "123456"
	now := time.Now().UTC()
	tm := now.Add(-1 * time.Duration(32) * time.Second)

	userDto := &dto.User{
		ID:          1,
		PhoneNumber: phoneNumber,
		Status:      constants.UserInitStatus,
		CreatedAt:   tm,
		UpdatedAt:   tm,
	}

	userStore := mockStores.NewMockIUserStore(ctrl)
	userStore.EXPECT().GetByPhoneNumber(gomock.Eq(phoneNumber)).Return(userDto, true, nil)
	expectedError := errors.New("Too many request ")
	userStore.EXPECT().UpdateStatus(gomock.Any()).Return(expectedError)

	userOtpStore := mockStores.NewMockIUserOtpStore(ctrl)
	userOtpDto := dto.UserOtp{
		ID:        2,
		UserID:    1,
		Otp:       "123456",
		CreatedAt: tm,
		UpdatedAt: tm,
	}

	userOtpStore.EXPECT().GetByUserID(gomock.Eq(userDto.ID)).Return(userOtpDto, true, nil)
	smsService := mockExternal.NewMockISmsService(ctrl)
	userValidator := validator.NewUserValidator()
	userOtpValidator := validator.NewUserOtpValidator()
	userOtpHelper := helpers.NewUserOtpHelper()
	userHelper := helpers.NewUserHelper("abc")

	cfg := config.Config{
		Base:                 config.Base{},
		MySQL:                config.MySQL{},
		PhoneNumberRateLimit: config.PhoneNumberRateLimit{},
		Otp:                  config.Otp{},
		SmsService:           config.SmsService{},
		Token:                config.Token{},
	}

	cfg.Otp.ExpiredTime = 60
	cfg.Otp.Size = 6
	userService := services.NewUserService(
		cfg,
		smsService,
		userValidator,
		userOtpValidator,
		userOtpHelper,
		userHelper,
		userStore,
		userOtpStore,
	)

	_, err := userService.Login(phoneNumber, otp)
	if err == nil || err.Error() != expectedError.Error() {
		t.Fatalf("expected error %v", expectedError)
	}
}

func TestUserService_Login_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	phoneNumber := "0961234567"
	otp := "123456"
	now := time.Now().UTC()
	tm := now.Add(-1 * time.Duration(32) * time.Second)

	userDto := &dto.User{
		ID:          1,
		PhoneNumber: phoneNumber,
		Status:      constants.UserInitStatus,
		CreatedAt:   tm,
		UpdatedAt:   tm,
	}

	userStore := mockStores.NewMockIUserStore(ctrl)
	userStore.EXPECT().GetByPhoneNumber(gomock.Eq(phoneNumber)).Return(userDto, true, nil)
	userStore.EXPECT().UpdateStatus(gomock.Any()).Return(nil)

	userOtpStore := mockStores.NewMockIUserOtpStore(ctrl)
	userOtpDto := dto.UserOtp{
		ID:        2,
		UserID:    1,
		Otp:       "123456",
		CreatedAt: tm,
		UpdatedAt: tm,
	}

	userOtpStore.EXPECT().GetByUserID(gomock.Eq(userDto.ID)).Return(userOtpDto, true, nil)
	smsService := mockExternal.NewMockISmsService(ctrl)
	userValidator := validator.NewUserValidator()
	userOtpValidator := validator.NewUserOtpValidator()
	userOtpHelper := helpers.NewUserOtpHelper()
	userHelper := helpers.NewUserHelper("abc")

	cfg := config.Config{
		Base:                 config.Base{},
		MySQL:                config.MySQL{},
		PhoneNumberRateLimit: config.PhoneNumberRateLimit{},
		Otp:                  config.Otp{},
		SmsService:           config.SmsService{},
		Token:                config.Token{},
	}

	cfg.Otp.ExpiredTime = 60
	cfg.Otp.Size = 6
	userService := services.NewUserService(
		cfg,
		smsService,
		userValidator,
		userOtpValidator,
		userOtpHelper,
		userHelper,
		userStore,
		userOtpStore,
	)

	token, err := userService.Login(phoneNumber, otp)
	if err != nil {
		t.Fatalf("expected nil")
	}

	if token != "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjoxfQ.6dJuEd4gOyQ3aKNHDABXEwQrTYQgNPVDLd0ZcyIUw-Q" {
		t.Fatalf("wrong token")
	}
}
