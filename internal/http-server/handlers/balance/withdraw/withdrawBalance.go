package withdraw

import (
	balanceService "accumulativeSystem/internal/services/balance"
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
		userId := r.Context().Value("user_id")

		if userId == nil {
			http.Error(w, "Not user id", http.StatusInternalServerError)
		}

		userID, ok := userId.(float64)

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

		if err.Error() == "insufficient funds" {
			http.Error(w, "there are insufficient funds in the account", http.StatusPaymentRequired)
			return
		}

		if err != nil {
			http.Error(w, "Error", http.StatusInternalServerError)
		}
	}
}
