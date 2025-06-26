package main

import (
	"encoding/json"
	"fmt"
	"math/rand/v2"
	"net/http"
	"os"
	"strconv"

	"github.com/NathanGdS/pkg"
)

func helloHandler(w http.ResponseWriter, r *http.Request) {
	var response pkg.Response
	if rand.Float32() < 0.5 {
		pkg.Error("Error on request :(")

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(500)
		response = pkg.Response{
			Message: "Error on request :(",
			Status:  500,
		}
		json.NewEncoder(w).Encode(response)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(201)
	response = pkg.Response{
		Message: "Success",
		Status:  201,
	}

	pkg.Default("Success; " + strconv.Itoa(response.Status) + " - " + response.Message)
	json.NewEncoder(w).Encode(response)
}

func main() {
	http.HandleFunc("/", helloHandler)
	err := http.ListenAndServe(":8081", nil)

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	} else {
		fmt.Println("Server listening on port 8081")
	}
}
