package rmqevents

const OrderCreatedEventExchange = "order.created"

type OrderCreatedEvent struct {
	OrderID uint64 `json:"order_id"`
}

func NewOrderCreatedEvent(orderId uint64) OrderCreatedEvent {
	return OrderCreatedEvent{OrderID: orderId}
}
