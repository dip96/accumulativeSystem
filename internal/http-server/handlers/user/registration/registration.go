package registration

import (
	//serviceUser "accumulativeSystem/internal/service/user" //TODO добавить отдельный слой service, прослойка между контроллерами и моделями
	userService "accumulativeSystem/internal/services/user"
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
		//TODO другой способ
		var req Request
		err := render.DecodeJSON(r.Body, &req)

		user, err := userService.CreateUser(req.Login, req.Password)

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
