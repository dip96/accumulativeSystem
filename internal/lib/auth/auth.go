package auth

import (
	userModel "accumulativeSystem/internal/models/user"
	"golang.org/x/crypto/bcrypt"
)

func Authenticate(user *userModel.User, password string) error {
	err := bcrypt.CompareHashAndPassword(user.HashPassword, []byte(password))
	if err != nil {
		return err
	}

	return nil
}
