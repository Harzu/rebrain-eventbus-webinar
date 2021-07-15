package kafkaevents

const PaymentSuccessTopic = "payment_success"

type PaymentSuccessEvent struct {
	UserID uint64 `json:"user_id"`
	Amount uint64 `json:"amount"`
}

func NewPaymentSuccessEvent(userId, amount uint64) PaymentSuccessEvent {
	return PaymentSuccessEvent{
		UserID: userId,
		Amount: amount,
	}
}
