package main

import (
	"errors"
	"fmt"
	"time"

	"github.com/sony/gobreaker"
)

var totalRequest = 0
var successes = 0
var failures = 0

func unstableAPI() (string, error) {
	totalRequest++
	fmt.Printf("Hit on API %d ", totalRequest)

	if totalRequest >= 15 {
		successes++
		return "success", nil
	}

	if totalRequest >= 5 {
		failures++
		return "", errors.New("simulated failure")
	}
	successes++
	return "success", nil
}

func main() {
	settings := gobreaker.Settings{
		Name:        "HTTP API",
		MaxRequests: 3,
		Interval:    0,               // intervalo para reset de contagem (0 = desativado)
		Timeout:     5 * time.Second, // tempo que o circuito fica aberto antes de testar de novo
		ReadyToTrip: func(counts gobreaker.Counts) bool {
			return counts.ConsecutiveFailures > 3
		},
		OnStateChange: func(name string, from gobreaker.State, to gobreaker.State) {
			fmt.Printf("Circuit breaker mudou de %v para %v\n", from, to)
		},
	}

	cb := gobreaker.NewCircuitBreaker(settings)

	for totalRequest < 20 {
		time.Sleep(500 * time.Millisecond)
		result, err := cb.Execute(func() (interface{}, error) {
			return unstableAPI()
		})

		if err != nil {
			fmt.Printf("Erro na requisição: %v\n", err)
		} else {
			fmt.Printf("Resposta: %v\n", result)
		}
	}

	fmt.Printf("\nTotal hits: %d\nTotal successes: %d\nTotal failures: %d\n ", totalRequest, successes, failures)
}
