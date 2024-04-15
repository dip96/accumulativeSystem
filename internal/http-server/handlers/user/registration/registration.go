package registration

import (
	errPostgres "accumulativeSystem/internal/errors/postgres"
	"accumulativeSystem/internal/lib/api/response"
	"accumulativeSystem/internal/lib/hash"
	userModel "accumulativeSystem/internal/model/user"
	storage "accumulativeSystem/internal/storage/postgres"
	"errors"
	_ "github.com/go-chi/jwtauth"
	"github.com/go-chi/render"
	"net/http"
)

type Request struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

type Response struct {
	response.Response
}

func New(postgres *storage.Postgres) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		user := &userModel.User{}

		// Связываем данные запроса с моделью User
		if err := render.Bind(r, user); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		hashPassword, err := hash.HashPassword(user.Password)

		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		user.HashPassword = hashPassword

		err = postgres.CreateUser(user)

		//TODO стоит ли в хендлере оставить это условие. Возможно стоит перенести в другой слой
		if err != nil {
			var postgresErr *errPostgres.PostgresError
			if errors.As(err, &postgresErr) {
				http.Error(w, postgresErr.Error(), http.StatusConflict)
			} else {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
		}

	}
}
