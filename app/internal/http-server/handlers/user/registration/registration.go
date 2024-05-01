package registration

import (
	apiError "accumulativeSystem/internal/errors/api"
	userService "accumulativeSystem/internal/services/user"
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

func New(userService userService.UserService, jwtAuth *jwtauth.JWTAuth) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req Request
		err := render.DecodeJSON(r.Body, &req)

		user, err := userService.CreateUser(req.Login, req.Password)

		if err != nil {
			var customErr *apiError.ApiError
			if errors.As(err, &customErr) {
				http.Error(w, customErr.Error(), customErr.Code)
				return
			}
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		_, tokenString, err := jwtAuth.Encode(map[string]interface{}{
			"user_id": user.Id,
			"exp":     time.Now().Add(time.Hour * 24).Unix(), // действителен в течение 24 часов
		})

		if err != nil {
			//TODO Нужно ли удалять пользователя? Или попросить сделать реавторизацию?
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}

		w.Header().Set("Authorization", tokenString)
		w.WriteHeader(http.StatusOK)
	}
}
