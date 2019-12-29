package services

import (
	"errors"
	"fmt"
	"log"
	"tbox_backend/config"
	"tbox_backend/external"
	"tbox_backend/internal/constants"
	"tbox_backend/internal/dto"
	e "tbox_backend/internal/errors"
	"tbox_backend/internal/helpers"
	"tbox_backend/internal/stores"
	"tbox_backend/internal/validator"
	"time"
)

type IUserService interface {
	GenerateOtp(phoneNumber string) error
	ResendOtp(phoneNumber string) error
	Login(phoneNumber string, otp string) (string, error)
}

type UserService struct {
	cfg              config.Config
	smsService       external.ISmsService
	userValidator    validator.IUserValidator
	userOtpValidator validator.IUserOtpValidator
	userOtpCommon    helpers.IUserOtpHelper
	userCommon       helpers.IUserHelper
	userStore        stores.IUserStore
	userOtpStore     stores.IUserOtpStore
}

func NewUserService(
	cfg config.Config,
	smsService external.ISmsService,
	userValidator validator.IUserValidator,
	userOtpValidator validator.IUserOtpValidator,
	userOtpCommon helpers.IUserOtpHelper,
	userCommon helpers.IUserHelper,
	userStore stores.IUserStore,
	userOtpStore stores.IUserOtpStore,
) *UserService {
	return &UserService{
		cfg:              cfg,
		smsService:       smsService,
		userValidator:    userValidator,
		userOtpValidator: userOtpValidator,
		userOtpCommon:    userOtpCommon,
		userCommon:       userCommon,
		userStore:        userStore,
		userOtpStore:     userOtpStore,
	}
}

func (s UserService) GenerateOtp(phoneNumber string) error {
	user, exists, err := s.userStore.GetByPhoneNumber(phoneNumber)
	if err != nil {
		return err
	}

	if !exists {
		user = &dto.User{
			PhoneNumber: phoneNumber,
			Status:      constants.UserInitStatus,
			CreatedAt:   time.Now().UTC(),
			UpdatedAt:   time.Now().UTC(),
		}

		err := s.userStore.Save(user)
		if err != nil {
			return err
		}
	} else if user.Status == constants.UserVerifiedStatus {
		return e.VerifiedPhoneNumberError{PhoneNumber: phoneNumber}
	}

	userOtp, exists, err := s.userOtpStore.GetByUserID(user.ID)
	if err != nil {
		return err
	}

	if exists {
		now := time.Now().UTC()
		if now.Sub(userOtp.UpdatedAt).Seconds() > float64(s.cfg.Otp.ExpiredTime) {
			otp := s.userOtpCommon.GenerateRandomOtp(s.cfg.Otp.Size)
			userOtp.Otp = otp
			userOtp.UpdatedAt = time.Now().UTC()
			err := s.userOtpStore.UpdateOtp(userOtp)
			if err != nil {
				return err
			}

			err = s.smsService.SendOtp(phoneNumber, otp)
			if err != nil {
				log.Println(fmt.Sprintf("Failed to send sms to %s", phoneNumber), err)
			}

			return nil
		} else {
			return e.GeneratedOtpError{}
		}
	} else {
		otp := s.userOtpCommon.GenerateRandomOtp(s.cfg.Otp.Size)
		err := s.userOtpStore.Save(dto.UserOtp{
			UserID:    user.ID,
			Otp:       otp,
			CreatedAt: time.Now().UTC(),
			UpdatedAt: time.Now().UTC(),
		})

		if err != nil {
			return err
		}

		err = s.smsService.SendOtp(phoneNumber, otp)
		if err != nil {
			log.Println(fmt.Sprintf("Failed to send sms to %s", phoneNumber), err)
		}

		return nil
	}
}

func (s UserService) ResendOtp(phoneNumber string) error {
	user, exists, err := s.userStore.GetByPhoneNumber(phoneNumber)
	if err != nil {
		return err
	} else if !exists {
		return e.NotExistsPhoneNumberError{PhoneNumber: phoneNumber}
	} else if user.Status == constants.UserVerifiedStatus {
		return e.VerifiedPhoneNumberError{PhoneNumber: phoneNumber}
	}

	userOtp, exists, err := s.userOtpStore.GetByUserID(user.ID)
	if err != nil {
		return err
	} else if !exists {
		return errors.New("Could not resend OTP ")
	}

	now := time.Now().UTC()
	if now.Sub(userOtp.UpdatedAt).Seconds() > float64(s.cfg.Otp.ResendWaitingTime) {
		otp := s.userOtpCommon.GenerateRandomOtp(s.cfg.Otp.Size)
		userOtp.Otp = otp
		userOtp.UpdatedAt = time.Now().UTC()
		err := s.userOtpStore.UpdateOtp(userOtp)
		if err != nil {
			return err
		}

		err = s.smsService.SendOtp(phoneNumber, otp)
		if err != nil {
			log.Println(fmt.Sprintf("Failed to send sms to %s", phoneNumber), err)
		}

		return nil
	} else {
		return e.GeneratedOtpError{}
	}
}

func (s UserService) Login(phoneNumber string, otp string) (string, error) {
	if valid := s.userValidator.IsPhoneNumberValid(phoneNumber); !valid {
		return "", e.InvalidPhoneNumberError{PhoneNumber: phoneNumber}
	}

	user, exists, err := s.userStore.GetByPhoneNumber(phoneNumber)
	if err != nil {
		return "", err
	} else if !exists {
		return "", e.NotExistsPhoneNumberError{PhoneNumber: phoneNumber}
	} else if user.Status == constants.UserVerifiedStatus {
		token, _ := s.userCommon.GenerateToken(user.ID)
		return token, nil
	}

	if valid := s.userOtpValidator.IsOtpValid(otp, s.cfg.Otp.Size); !valid {
		return "", e.InvalidOtpError{Otp: otp}
	}

	userOtp, exists, err := s.userOtpStore.GetByUserID(user.ID)
	if err != nil {
		return "", err
	} else if !exists {
		return "", e.IncorrectOtpError{Otp: otp}
	}

	now := time.Now().UTC()
	if now.Sub(userOtp.UpdatedAt).Seconds() <= float64(s.cfg.Otp.ExpiredTime) {
		if otp == userOtp.Otp {
			user.Status = constants.UserVerifiedStatus
			user.UpdatedAt = time.Now().UTC()
			err := s.userStore.UpdateStatus(user)
			if err != nil {
				return "", err
			}

			token, _ := s.userCommon.GenerateToken(user.ID)
			return token, nil
		} else {
			return "", e.IncorrectOtpError{Otp: otp}
		}
	} else {
		return "", e.ExpiredOtpError{Otp: otp}
	}
}
