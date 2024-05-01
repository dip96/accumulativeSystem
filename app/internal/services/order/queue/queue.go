package orderChan

type OrderQueueService interface {
	EnqueueOrder(orderId int)
	GetOrderChan() chan int
}

type orderQueueService struct {
	orderChan chan int
}

func NewOrderQueue() OrderQueueService {
	orderChan := make(chan int)
	return &orderQueueService{
		orderChan: orderChan,
	}
}

func (s *orderQueueService) EnqueueOrder(orderID int) {
	s.orderChan <- orderID
}

func (s *orderQueueService) GetOrderChan() chan int {
	return s.orderChan
}
