package rabbitmq

import "time"

type Strategy interface {
	NextBackoff() time.Duration
}

type linearBackoff time.Duration

func NewLinearBackoff(backoff time.Duration) Strategy {
	return linearBackoff(backoff)
}

func (r linearBackoff) NextBackoff() time.Duration {
	return time.Duration(r)
}
