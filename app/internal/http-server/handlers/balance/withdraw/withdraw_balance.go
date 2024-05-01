package withdraw

import (
	APIError "accumulativeSystem/internal/errors/api"
	balanceService "accumulativeSystem/internal/services/balance"
	"errors"
	"github.com/go-chi/render"
	"net/http"
	"strconv"
)

type WithdrawBalanceRequest struct {
	Order string  `json:"order"`
	Sum   float64 `json:"sum"`
}

func New(service balanceService.BalanceService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		contextUserID := r.Context().Value("user_id")

		if contextUserID == nil {
			http.Error(w, "Not user id", http.StatusInternalServerError)
		}

		userID, ok := contextUserID.(int)

		if !ok {
			http.Error(w, "Error userID", http.StatusInternalServerError)
			return
		}

		var req WithdrawBalanceRequest
		err := render.DecodeJSON(r.Body, &req)

		if err != nil {
			http.Error(w, "Error", http.StatusInternalServerError)
			return
		}

		orderIDInt, err := strconv.Atoi(req.Order)

		if err != nil {
			http.Error(w, "Error", http.StatusInternalServerError)
			return
		}

		err = service.WithdrawBalance(int(userID), orderIDInt, req.Sum)

		if err != nil {
			var customErr *APIError.APIError
			if errors.As(err, &customErr) {
				http.Error(w, customErr.Error(), customErr.Code)
				return
			}
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
	}
}
