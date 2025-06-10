package main

import (
	"fmt"
	"math/rand/v2"
	"net/http"
	"os"

	"github.com/NathanGdS/pkg"
)

func helloHandler(w http.ResponseWriter, r *http.Request) {
	if rand.Float32() < 0.1 {
		pkg.Error("Error on request :(")
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(500)
		return
	}

	pkg.Default("Success")
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(201)
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
