package auth

import (
	"accumulativeSystem/internal/lib/hash"
	userModel "accumulativeSystem/internal/models/user"
	"golang.org/x/crypto/bcrypt"
)

func Authenticate(user *userModel.User, password string) error {
	//TODO стоит ли в хендлере оставить это условие. Возможно стоит перенести в другой слой
	hashPassword, err := hash.HashPassword(password)
	if err != nil {
		return err
	}
	//END TODO

	// Сравнить хэш пароля из базы данных с хэшем пароля, полученным в запросе
	if err := bcrypt.CompareHashAndPassword(user.HashPassword, hashPassword); err != nil {
		return err
	}

	return nil
}
