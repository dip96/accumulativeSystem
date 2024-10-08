package registration

import (
	APIError "accumulativeSystem/internal/errors/api"
	"accumulativeSystem/internal/http-server/middleware/auth"
	orderService "accumulativeSystem/internal/services/order"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"time"
)

type Request struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

type OrderResponse struct {
	Order      string    `json:"order"`
	Sum        float32   `json:"sum"`
	UploadedAt time.Time `json:"processed_at"`
}

func New(service orderService.OrderService) http.HandlerFunc {
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

		orders, err := service.GetWithdrawalsByUserID(userID)

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
				Order:      strconv.Itoa(order.OrderID),
				Sum:        float32(order.WithdrawnBalance),
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
