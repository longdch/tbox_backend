package helpers

import (
	"github.com/dgrijalva/jwt-go"
)

type IUserHelper interface {
	GenerateToken(userID int) (string, error)
}

type UserHelper struct {
	secretKey string
}

func NewUserHelper(secretKey string) *UserHelper {
	return &UserHelper{secretKey: secretKey}
}

func (c UserHelper) GenerateToken(userID int) (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	claims["user_id"] = userID
	return token.SignedString([]byte(c.secretKey))
}
