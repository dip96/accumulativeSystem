package create

import (
	APIError "accumulativeSystem/internal/errors/api"
	"accumulativeSystem/internal/http-server/middleware/auth"
	orderModel "accumulativeSystem/internal/models/order"
	orderService "accumulativeSystem/internal/services/order"
	"errors"
	"io"
	"net/http"
	"strconv"
	"strings"
)

func New(service orderService.OrderService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var userIDKey = auth.ContextKey("user_id")
		contextUserID := r.Context().Value(userIDKey)

		if contextUserID == nil {
			http.Error(w, "Not user id", http.StatusInternalServerError)
		}

		userID, ok := contextUserID.(int)

		if !ok {
			http.Error(w, "Error", http.StatusInternalServerError)
			return
		}

		body, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		orderID := strings.TrimSpace(string(body))

		orderIDInt, err := strconv.Atoi(orderID)

		if err != nil {
			http.Error(w, "Error userID", http.StatusInternalServerError)
			return
		}

		var order orderModel.Order
		order.OrderID = orderIDInt
		order.UserID = userID

		_, err = service.CreateOrder(&order)

		if err != nil {
			var customErr *APIError.APIError
			if errors.As(err, &customErr) {
				http.Error(w, customErr.Error(), customErr.Code)
				return
			}
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusAccepted)
	}
}
