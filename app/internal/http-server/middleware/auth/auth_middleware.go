package auth

import (
	"accumulativeSystem/internal/logger"
	"context"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/jwtauth"
	"github.com/lestrrat-go/jwx/jwt"
	"net/http"
	"time"
)

type ContextKey string

func AuthMiddleware(log logger.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()
			defer func() {
				duration := time.Since(start)
				requestID := middleware.GetReqID(r.Context())
				log.Infof("Request ID: %s, Request: %s %s took %v", requestID, r.Method, r.URL.Path, duration)
			}()

			//TODo в случаи удаления пользователя токен корректно отрабатывает
			tokenStr := jwtauth.TokenFromHeader(r)

			if tokenStr == "" {
				tokenStr = r.Header.Get("Authorization")
			}

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

			var userIDKey = ContextKey("user_id")
			userID := int(mapValue["user_id"].(float64))
			ctx := context.WithValue(r.Context(), userIDKey, userID)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
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
