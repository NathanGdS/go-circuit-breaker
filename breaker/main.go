package main

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"strconv"
	"time"

	"github.com/NathanGdS/pkg"
	"github.com/sony/gobreaker"
)

var totalRequest = 0
var successes = 0
var failures = 0

func unstableAPI() (pkg.Response, error) {
	client := &http.Client{}
	req, err := http.NewRequest("GET", "http://localhost:8081", nil)
	if err != nil {
		pkg.Error("Error creating request:  " + err.Error())
		return pkg.Response{}, errors.New("failure")
	}

	resp, err := client.Do(req)
	if err != nil {
		pkg.Error("Error sending request:  " + err.Error())
		return pkg.Response{}, errors.New("failure")
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		pkg.Error("Error reading response body: " + err.Error())
		return pkg.Response{}, errors.New("failure")
	}

	var response pkg.Response
	err = json.Unmarshal(body, &response)
	if err != nil {
		pkg.Error("Error unmarshalling response: " + err.Error())
		return pkg.Response{}, errors.New("failure")
	}

	if response.Status >= 400 {
		return pkg.Response{}, errors.New(response.Message)

	}

	return response, nil
}

func executeWithCircuitBreaker(cb *gobreaker.CircuitBreaker) (pkg.Response, error) {
	response, err := cb.Execute(func() (interface{}, error) {
		return unstableAPI()
	})

	if err != nil {
		return pkg.Response{}, err
	}

	if response, ok := response.(pkg.Response); ok {
		return response, nil
	}

	return pkg.Response{}, errors.New("unexpected response type")
}

func main() {
	settings := gobreaker.Settings{
		Name:        "HTTP API",
		MaxRequests: 3,
		Interval:    0,               // intervalo para reset de contagem (0 = desativado)
		Timeout:     5 * time.Second, // tempo que o circuito fica aberto antes de testar de novo
		ReadyToTrip: func(counts gobreaker.Counts) bool {
			return counts.ConsecutiveFailures >= 3
		},
		OnStateChange: func(name string, from gobreaker.State, to gobreaker.State) {
			pkg.Default("Circuit breaker mudou de " + from.String() + " para " + to.String())
		},
	}

	cb := gobreaker.NewCircuitBreaker(settings)

	for {
		time.Sleep(50 * time.Millisecond)
		response, err := executeWithCircuitBreaker(cb)

		if err != nil {
			pkg.Error("Erro na requisição: " + err.Error())
		} else {
			pkg.Default("Resposta: " + response.Message + " - " + strconv.Itoa(response.Status))
		}
	}
}
