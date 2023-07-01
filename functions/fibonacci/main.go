package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strconv"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type response struct {
	Result      int    `json:"result"`
	ProcessedBy string `json:"processed_by"`
}

var (
	requests = prometheus.NewCounter(prometheus.CounterOpts{
		Name:        "success_requests_total",
		Help:        "Number of all successful requests.",
		ConstLabels: prometheus.Labels{"instance": "function/fibonacci", "path": "success_requests_total", "method": "GET"},
	})
)

func init() {
	prometheus.MustRegister(requests)
}

func main() {
	fmt.Println("Executing fibonacci function.")

	http.HandleFunc("/", fibonacciHandler)
	http.Handle("/metrics", promhttp.Handler())

	err := http.ListenAndServe(":3333", nil)
	if err != nil {
		fmt.Println("Error occurred running server: ", err)
	}
}

// Handler for fibonacci calculator
func fibonacciHandler(w http.ResponseWriter, r *http.Request) {
	number, err := strconv.Atoi(r.URL.Query().Get("number"))
	fmt.Println("Got a request: ", number)
	if err != nil {
		fmt.Println("Wrong input!")
	}

	result := fibonacci(number)

	response := response{
		result,
		os.Getenv("HOSTNAME"),
	}

	responseBytes, err := json.Marshal(response)
	if err != nil {
		fmt.Println("500 internal error!")
	}

	requests.Inc()
	w.Header().Set("Content-Type", "application/json")
	w.Write(responseBytes)
}

// Returns the fibonacci of a given number
func fibonacci(n int) int {
	if n < 2 {
		return n
	}

	return fibonacci(n-1) + fibonacci(n-2)
}
