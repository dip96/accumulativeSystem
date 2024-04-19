package getBalance

import (
	//serviceUser "accumulativeSystem/internal/service/user" //TODO добавить отдельный слой service, прослойка между контроллерами и моделями
	storage "accumulativeSystem/internal/storage/postgres"
	"encoding/json"
	"net/http"
)

type BalanceResponse struct {
	Balance          float64 `json:"current"`
	WithdrawnBalance float64 `json:"withdrawn"`
}

func New(postgres *storage.Postgres) http.HandlerFunc {
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

		balance, err := postgres.GetUserBalance((int(userID)))

		if err != nil {
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
