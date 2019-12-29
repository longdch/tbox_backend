package stores

import (
	"database/sql"
	"errors"
	"github.com/jmoiron/sqlx"
	"tbox_backend/internal/dto"
	"tbox_backend/internal/models"
)

type IUserStore interface {
	GetByPhoneNumber(phoneNo string) (*dto.User, bool, error)
	Save(user *dto.User) error
	UpdateStatus(user *dto.User) error
}

type UserStore struct {
	client *sqlx.DB
}

func NewUserStore(client *sqlx.DB) *UserStore {
	return &UserStore{client: client}
}

func (s *UserStore) GetByPhoneNumber(phoneNo string) (*dto.User, bool, error) {
	query := `
	SELECT u.user_id,
	u.phone_number,
	u.status,
	u.created_at,
	u.updated_at
	FROM users u
	WHERE u.phone_number = ?
	`

	userModel := models.User{}
	err := s.client.Get(&userModel, query, phoneNo)
	if err != nil && err == sql.ErrNoRows {
		return nil, false, nil
	} else if err != nil {
		return nil, true, err
	} else {
		userDto := userModel.ToDto()
		return &userDto, true, nil
	}
}

func (s *UserStore) Save(user *dto.User) error {
	query := `
	INSERT INTO users (user_id, phone_number, status, created_at, updated_at) 
	VALUES (:user_id, :phone_number, :status, :created_at, :updated_at)
	`

	userModel := &models.User{}
	userModel.FromDto(user)
	insertResult, err := s.client.NamedExec(query, &userModel)
	if err != nil {
		return err
	}

	if insertResult != nil {
		lastInsertID, err := insertResult.LastInsertId()
		if err != nil {
			return err
		} else {
			user.ID = int(lastInsertID)
			return nil
		}
	} else {
		return errors.New("Could not insert user ")
	}
}

func (s *UserStore) UpdateStatus(user *dto.User) error {
	query := `
	UPDATE users SET status = :status, updated_at = :updated_at WHERE user_id = :user_id
	`

	userModel := &models.User{}
	userModel.FromDto(user)
	_, err := s.client.NamedExec(query, userModel)
	return err
}
