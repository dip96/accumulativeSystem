package get

import (
	APIError "accumulativeSystem/internal/errors/api"
	orderService "accumulativeSystem/internal/services/order"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"time"
)

type Request struct {
	orderID string
}

type OrderResponse struct {
	Number     string    `json:"number"`
	Status     string    `json:"status"`
	Accrual    float32   `json:"accrual"`
	UploadedAt time.Time `json:"uploaded_at"`
}

func New(service orderService.OrderService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		contextUserID := r.Context().Value("user_id")

		if contextUserID == nil {
			http.Error(w, "Not user id", http.StatusInternalServerError)
		}

		userID, ok := contextUserID.(int)

		if !ok {
			http.Error(w, "Error userID", http.StatusInternalServerError)
		}

		orders, err := service.GetOrdersByUserID(userID)

		if err != nil {
			var customErr *APIError.APIError
			if errors.As(err, &customErr) {
				http.Error(w, customErr.Error(), customErr.Code)
				return
			}
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		orderResponses := make([]OrderResponse, len(orders))
		for i, order := range orders {
			createdAt := order.CreatedAt.Time
			orderResponses[i] = OrderResponse{
				Number:     strconv.Itoa(order.OrderID),
				Status:     string(order.Status),
				Accrual:    float32(order.Accrual),
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
