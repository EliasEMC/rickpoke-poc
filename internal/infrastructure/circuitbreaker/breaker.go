package circuitbreaker

import (
	"time"

	"github.com/sony/gobreaker"
)

func New(name string) *gobreaker.CircuitBreaker {
	settings := gobreaker.Settings{
		Name:        name,
		MaxRequests: 3,
		Timeout:     5 * time.Second,
		ReadyToTrip: func(c gobreaker.Counts) bool {
			return c.ConsecutiveFailures >= 5
		},
	}
	return gobreaker.NewCircuitBreaker(settings)
}
