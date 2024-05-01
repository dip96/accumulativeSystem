package registration

import (
	"accumulativeSystem/internal/lib/auth"
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

		if err != nil {
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}

		user, err := userService.GetUserWithPassword(req.Login)

		if err != nil {
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}

		err = auth.Authenticate(user, req.Password)

		if err != nil {
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}

		_, tokenString, err := jwtAuth.Encode(map[string]interface{}{
			"user_id": user.ID,
			"exp":     time.Now().Add(time.Hour * 24).Unix(), // действителен в течение 24 часов
		})

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}

		w.Header().Set("Authorization", tokenString)
		w.WriteHeader(http.StatusOK)
	}
}
