package createOrder

import (
	storage "accumulativeSystem/internal/storage/postgres"
	"io"
	"net/http"
	"strconv"
	"strings"
)

type Request struct {
	orderID string
}

func New(postgres *storage.Postgres) http.HandlerFunc {
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

		//var req Request
		//req.orderID = orderID

		if !isValidLunaChecksum(orderID) {
			http.Error(w, "Invalid order ID", http.StatusUnprocessableEntity)
			return
		}

		orderIDInt, err := strconv.Atoi(orderID)
		userID, ok := userId.(float64)

		if !ok {
			http.Error(w, "Error userID", http.StatusInternalServerError)
		}

		existingOrder, err := postgres.GetOrderByOrderId(orderIDInt)
		if existingOrder.UserId != userId {
			http.Error(w, "Order already exists for another user", http.StatusConflict)
			return
		}

		if err == nil && existingOrder != nil {
			http.Error(w, "Order already exists", http.StatusOK)
			return
		}

		postgres.CreateOrder(userID, orderIDInt)
		//order, err := postgres.CreateOrder(userID, orderIDInt)

		if err != nil {
			http.Error(w, "Error order", http.StatusInternalServerError)
		}

		w.WriteHeader(http.StatusAccepted)
	}
}

func isValidLunaChecksum(creditCardNumber string) bool {
	var sum int
	var isEven = false

	for i := len(creditCardNumber) - 1; i >= 0; i-- {
		digit, _ := strconv.Atoi(string(creditCardNumber[i]))
		if isEven {
			digit *= 2
			if digit > 9 {
				digit -= 9
			}
		}
		sum += digit
		isEven = !isEven
	}

	return sum%10 == 0
}
