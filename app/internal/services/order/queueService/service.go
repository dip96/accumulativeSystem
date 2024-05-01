package queueservice

import (
	"accumulativeSystem/internal/config"
	orderModel "accumulativeSystem/internal/models/order"
	"accumulativeSystem/internal/services/balance"
	orderService "accumulativeSystem/internal/services/order"
	orderChan "accumulativeSystem/internal/services/order/queue"
	"encoding/json"
	"log"
	"net/http"
	"strconv"
)

type OrderQueue interface {
	RunGoroutine(service orderService.OrderService)
}

type orderQueue struct {
	urlAccrual string
	orderChan  orderChan.OrderQueueService
	usBalance  balance.BalanceService
}

type orderResponse struct {
	Order   string  `json:"order"`
	Status  string  `json:"status"`
	Accrual float64 `json:"accrual"`
}

func NewOrderQueueService(cfg config.ConfigInstance, service orderChan.OrderQueueService, balance balance.BalanceService) OrderQueue {
	return &orderQueue{
		urlAccrual: cfg.GetAccrualSystemAddress() + "/api/orders/",
		orderChan:  service,
		usBalance:  balance,
	}
}

func (s *orderQueue) RunGoroutine(service orderService.OrderService) {
	go func() {
		client := &http.Client{}
		for orderID := range s.orderChan.GetOrderChan() {

			order, err := service.GetOrderByOrderID(orderID)

			if err != nil {
				log.Printf("Problem searching for an order - %s", err.Error())
				continue
			}

			url := s.urlAccrual + strconv.Itoa(orderID)

			req, err := http.NewRequest("GET", url, nil)
			if err != nil {
				log.Printf("Error creating request: %v", err)
				continue
			}

			resp, err := client.Do(req)
			if err != nil {
				log.Printf("Error sending request: %v", err)
				continue
			}

			if resp.StatusCode != http.StatusOK {
				log.Printf("Error: Received status code %d, method: %s", resp.StatusCode, url)
				continue
			}

			var orderRes orderResponse

			err = json.NewDecoder(resp.Body).Decode(&orderRes)
			if err != nil {
				log.Printf("Error decoding response: %v", err)
				continue
			}

			order.Status = orderModel.GetOrderStatusByValue(orderRes.Status)
			order.Accrual = orderRes.Accrual

			err = s.usBalance.AccrualBalance(order.UserID, order, order.Accrual)

			if err != nil {
				log.Printf("Error saving order: %v", err)
			}
		}
	}()
}
