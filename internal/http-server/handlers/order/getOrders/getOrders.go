package createOrder

import (
	storage "accumulativeSystem/internal/storage/postgres"
	"encoding/json"
	"net/http"
	"time"
)

type Request struct {
	orderID string
}

type OrderResponse struct {
	Number     int64     `json:"number"`
	Status     string    `json:"status"`
	Accrual    float64   `json:"accrual"`
	UploadedAt time.Time `json:"uploaded_at"`
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

		orders, err := postgres.GetOrdersByUserId(int(userID))
		if len(orders) == 0 {
			http.Error(w, "", http.StatusNoContent)
			return
		}

		if err != nil {
			http.Error(w, "Order already exists", http.StatusInternalServerError)
			return
		}

		orderResponses := make([]OrderResponse, len(orders))
		for i, order := range orders {
			createdAt := order.CreatedAt.Time
			orderResponses[i] = OrderResponse{
				Number:     order.OrderId,
				Status:     string(order.Status),
				Accrual:    order.Accrual,
				UploadedAt: createdAt,
			}
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		jsonData, err := json.Marshal(orderResponses)

		if err != nil {
			http.Error(w, "Order already exists", http.StatusInternalServerError)
			return
		}

		_, err = w.Write(jsonData)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}
