package auth

import (
	"context"
	"github.com/go-chi/jwtauth"
	"github.com/lestrrat-go/jwx/jwt"
	"net/http"
)

func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		//TODo в случаи удаления пользователя токен корректно отрабатывает
		tokenStr := jwtauth.TokenFromHeader(r)
		if tokenStr == "" {
			http.Error(w, "Missing token", http.StatusUnauthorized)
			return
		}

		token, err := jwt.Parse([]byte(tokenStr))

		if err != nil {
			http.Error(w, "error token", http.StatusUnauthorized)
			return
		}

		mapValue, err := token.AsMap(r.Context())

		if err != nil {
			http.Error(w, "error token", http.StatusUnauthorized)
			return
		}

		if mapValue == nil {
			http.Error(w, "error token", http.StatusUnauthorized)
		}

		//TODO добавить проверку на актуальность токена по времени
		//expTime, ok := mapValue["exp"].(float64)
		//if !ok {
		//	http.Error(w, "Invalid token expiration", http.StatusUnauthorized)
		//	return
		//}

		userId := int(mapValue["user_id"].(float64))
		ctx := context.WithValue(r.Context(), "user_id", userId)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func JWTVerifier(jwt *jwtauth.JWTAuth) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			handler := jwtauth.Verifier(jwt)
			if handler != nil {
				// Обработка ошибки верификации токена
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}
