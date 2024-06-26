package get

import (
	APIError "accumulativeSystem/internal/errors/api"
	"accumulativeSystem/internal/http-server/middleware/auth"
	balanceService "accumulativeSystem/internal/services/balance"
	"encoding/json"
	"errors"
	"net/http"
)

type BalanceResponse struct {
	Balance          float64 `json:"current"`
	WithdrawnBalance float64 `json:"withdrawn"`
}

func New(service balanceService.BalanceService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var userIDKey = auth.ContextKey("user_id")
		contextUserID := r.Context().Value(userIDKey)

		if contextUserID == nil {
			http.Error(w, "Not user id", http.StatusInternalServerError)
		}

		userID, ok := contextUserID.(int)

		if !ok {
			http.Error(w, "Error userID", http.StatusInternalServerError)
		}

		balance, err := service.GetUserBalance(userID)

		if err != nil {
			var customErr *APIError.APIError
			if errors.As(err, &customErr) {
				http.Error(w, customErr.Error(), customErr.Code)
				return
			}
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		balRes := BalanceResponse{
			Balance:          balance.Balance,
			WithdrawnBalance: balance.WithdrawnBalance,
		}

		jsonData, err := json.Marshal(balRes)

		if err != nil {
			http.Error(w, "Order already exists", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		_, err = w.Write(jsonData)

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}
