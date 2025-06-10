package main

import (
	"errors"
	"net/http"
	"time"

	"github.com/NathanGdS/pkg"
	"github.com/sony/gobreaker"
)

var totalRequest = 0
var successes = 0
var failures = 0

func unstableAPI() (string, error) {
	client := &http.Client{}
	req, err := http.NewRequest("GET", "http://localhost:8081", nil)
	if err != nil {
		pkg.Error("Error creating request:  " + err.Error())
		return "", errors.New("failure")
	}

	resp, err := client.Do(req)
	if err != nil {
		pkg.Error("Error sending request:  " + err.Error())
		return "", errors.New("failure")
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		pkg.Error("Error on request!")
		return "", errors.New("failure")
	}

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
			pkg.Default("Circuit breaker mudou de " + from.String() + "para " + to.String())
		},
	}

	cb := gobreaker.NewCircuitBreaker(settings)

	for {
		// time.Sleep(500 * time.Millisecond)
		_, err := cb.Execute(func() (interface{}, error) {
			return unstableAPI()
		})

		if err != nil {
			pkg.Error("Erro na requisição: " + err.Error())
		}
	}

	// fmt.Printf("\nTotal hits: %d\nTotal successes: %d\nTotal failures: %d\n ", totalRequest, successes, failures)
}
