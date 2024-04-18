package registration

import (
	errPostgres "accumulativeSystem/internal/errors/postgres"
	"accumulativeSystem/internal/lib/hash"
	//serviceUser "accumulativeSystem/internal/service/user" //TODO добавить отдельный слой service, прослойка между контроллерами и моделями
	storage "accumulativeSystem/internal/storage/postgres"
	"errors"
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

		//TODO стоит ли в хендлере оставить это условие. Возможно стоит перенести в другой слой
		hashPassword, err := hash.HashPassword(req.Password)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		//END TODO

		user, err := postgres.CreateUser(req.Login, hashPassword)

		//TODO стоит ли в хендлере оставить это условие. Возможно стоит перенести в другой слой
		if err != nil {
			var postgresErr *errPostgres.PostgresError
			if errors.As(err, &postgresErr) {
				http.Error(w, postgresErr.Error(), http.StatusConflict)
			} else {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}

			return
		}
		//END TODO

		_, tokenString, err := jwtAuth.Encode(map[string]interface{}{
			"user_id": user.Id,
			"exp":     time.Now().Add(time.Hour * 24).Unix(), // действителен в течение 24 часов
		})

		if err != nil {
			//TODO Нужно ли удалять пользователя? Или попросить сделать реавторизацию?
		}

		w.Header().Set("Authorization", tokenString)
	}
}
