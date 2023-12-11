package circuit_breaker

import (
	"fmt"
	"time"

	"github.com/sony/gobreaker"
)

var CircuitBreaker *gobreaker.CircuitBreaker

func init() {
	CircuitBreaker = gobreaker.NewCircuitBreaker(gobreaker.Settings{
		Name:        "My Circuit Breaker",
		MaxRequests: 0,               // Maximum number of consecutive failures before tripping the circuit
		Interval:    5 * time.Second, // Duration to wait before allowing another request
		Timeout:     2 * time.Second, // Timeout for a single request attempt
		ReadyToTrip: func(counts gobreaker.Counts) bool {
			// failureRatio := float64(counts.TotalFailures) / float64(counts.Requests)
			// return counts.Requests >= 3 && failureRatio >= 0.6 // Trips the circuit if failure rate exceeds 60%
			return counts.ConsecutiveFailures > 3
		},
		OnStateChange: func(name string, from gobreaker.State, to gobreaker.State) {
			fmt.Printf("Circuit breaker '%s' changed from '%s' to '%s'\n", name, from, to)
		},
	})
}
