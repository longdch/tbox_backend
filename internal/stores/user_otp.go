package stores

import (
	"database/sql"
	"github.com/jmoiron/sqlx"
	"tbox_backend/internal/dto"
	"tbox_backend/internal/models"
)

type IUserOtpStore interface {
	GetByUserID(userID int) (dto.UserOtp, bool, error)
	Save(userOtp dto.UserOtp) error
	UpdateOtp(userOtp dto.UserOtp) error
}

type UserOtpStore struct {
	client *sqlx.DB
}

func NewUserOtpStore(client *sqlx.DB) *UserOtpStore {
	return &UserOtpStore{client: client}
}

func (s *UserOtpStore) GetByUserID(userID int) (dto.UserOtp, bool, error) {
	query := `
	SELECT u.user_otp_id,
	u.user_id,
	u.otp,
	u.created_at,
	u.updated_at
	FROM user_otp u
	WHERE u.user_id = ?
	`

	userOtpModel := models.UserOtp{}
	err := s.client.Get(&userOtpModel, query, userID)
	if err != nil && err == sql.ErrNoRows {
		return dto.UserOtp{}, false, nil
	} else if err != nil {
		return dto.UserOtp{}, false, err
	} else {
		return userOtpModel.ToOtp(), true, nil
	}
}

func (s *UserOtpStore) UpdateOtp(userOtp dto.UserOtp) error {
	query := `
	UPDATE user_otp SET otp = :otp, updated_at = :updated_at WHERE user_otp_id = :user_otp_id
	`

	userOtpModel := &models.UserOtp{}
	userOtpModel.FromDto(userOtp)
	_, err := s.client.NamedExec(query, userOtpModel)
	return err
}

func (s *UserOtpStore) Save(userOtp dto.UserOtp) error {
	query := `
	INSERT INTO user_otp (user_id, otp, created_at, updated_at) 
	VALUES (:user_id, :otp, :created_at, :updated_at)
	`

	userOtpModel := &models.UserOtp{}
	userOtpModel.FromDto(userOtp)
	_, err := s.client.NamedExec(query, &userOtpModel)
	return err
}
