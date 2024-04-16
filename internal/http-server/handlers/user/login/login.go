package registration

import (
	"accumulativeSystem/internal/lib/hash"
	userModel "accumulativeSystem/internal/models/user"
	"golang.org/x/crypto/bcrypt"

	//serviceUser "accumulativeSystem/internal/service/user"
	storage "accumulativeSystem/internal/storage/postgres"
	"github.com/go-chi/jwtauth"
	"github.com/go-chi/render"
	"net/http"
	"time"
)

type Request struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

func New(postgres *storage.Postgres, jwtAuth *jwtauth.JWTAuth) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		//TODO другой способ
		var req Request
		err := render.DecodeJSON(r.Body, &req)

		user, err := postgres.GetUserPassword(req.Login)

		if err != nil {
			//Ошибки - no rows in result set
			//Ошибки - context deadline exceeded

			//TODO возможно ошибка свянная с бд, а не с отсутствующим логином
			//http.Error(w, err.Error(), http.StatusInternalServerError)
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}

		err = Authenticate(user, req.Password)

		if err != nil {
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}

		_, tokenString, err := jwtAuth.Encode(map[string]interface{}{
			"user_id": user.Id,
			"exp":     time.Now().Add(time.Hour * 24).Unix(), // действителен в течение 24 часов
		})

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}

		w.Header().Set("Authorization", "Bearer "+tokenString)
	}
}

func Authenticate(user *userModel.User, password string) error {
	//TODO стоит ли в хендлере оставить это условие. Возможно стоит перенести в другой слой
	hashPassword, err := hash.HashPassword(password)
	if err != nil {
		return err
	}
	//END TODO

	// Сравнить хэш пароля из базы данных с хэшем пароля, полученным в запросе
	if err := bcrypt.CompareHashAndPassword(hashPassword, []byte(password)); err != nil {
		return err
	}

	return nil
}
