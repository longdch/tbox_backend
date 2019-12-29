package helpers

import (
	"math/rand"
	"time"
)

type IUserOtpHelper interface {
	GenerateRandomOtp(size int) string
}

type UserOtpHelper struct{}

func NewUserOtpHelper() *UserOtpHelper {
	return &UserOtpHelper{}
}

func (UserOtpHelper) GenerateRandomOtp(size int) string {
	var sources = [...]byte{
		'0', '1', '2', '3', '4', '5', '6', '7', '8', '9',
	}

	random := rand.New(rand.NewSource(time.Now().UnixNano()))
	otp := make([]byte, size)
	for i := 0; i < size; i++ {
		index := random.Int() % len(sources)
		otp[i] = sources[index]
	}

	return string(otp)
}
