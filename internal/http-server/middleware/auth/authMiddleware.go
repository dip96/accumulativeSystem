package auth

import (
	"context"
	"github.com/go-chi/jwtauth"
	"net/http"
)

func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token, claims, err := jwtauth.FromContext(r.Context())
		if err != nil || token == nil {
			http.Error(w, "Invalid token", http.StatusUnauthorized)
			return
		}

		if claims == nil {
			http.Error(w, "Invalid token claims", http.StatusUnauthorized)
			return
		}

		//TODO добавить проверку на актуальность токена
		//expDate := claims["exp"]

		ctx := context.WithValue(r.Context(), "user_id", claims["user_id"])
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
