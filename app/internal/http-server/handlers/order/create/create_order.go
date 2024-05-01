package create

import (
	apiError "accumulativeSystem/internal/errors/api"
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
		userId := r.Context().Value("user_id")

		if userId == nil {
			http.Error(w, "Not user id", http.StatusInternalServerError)
		}

		body, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		orderID := strings.TrimSpace(string(body))

		orderIDInt, err := strconv.Atoi(orderID)
		userID, ok := userId.(int)

		if !ok {
			http.Error(w, "Error userID", http.StatusInternalServerError)
		}

		var order orderModel.Order
		order.OrderId = orderIDInt
		order.UserId = userID

		_, err = service.CreateOrder(&order)

		if err != nil {
			var customErr *apiError.ApiError
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
