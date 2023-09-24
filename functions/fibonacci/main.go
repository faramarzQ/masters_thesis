package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"k8s.io/klog"
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

	timeout int
)

func init() {
	prometheus.MustRegister(requests)
}

func main() {
	fmt.Println("Executing fibonacci function.")

	http.HandleFunc("/", fibonacciHandler)
	http.Handle("/metrics", promhttp.Handler())

	var err error
	timeout, err = strconv.Atoi(os.Getenv("FIBONACCI_TIMEOUT_SECONDS"))
	if err != nil {
		klog.Error("Failed reading env var: FIBONACCI_TIMEOUT_SECONDS")
	}

	err = http.ListenAndServe(":3333", nil)
	if err != nil {
		klog.Fatal("Error occurred running server: ", err)
	}
}

// Handler for fibonacci calculator
func fibonacciHandler(w http.ResponseWriter, r *http.Request) {
	number, err := strconv.Atoi(r.URL.Query().Get("number"))
	klog.Info("Got a request: ", number)
	if err != nil {
		klog.Error("Wrong input!")
	}

	resChan := make(chan int)
	go calculateFibonacci(number, resChan)

	response := response{
		ProcessedBy: os.Getenv("HOSTNAME"),
	}

	select {
	case result := <-resChan:
		response.Result = result
	case <-time.After(time.Duration(timeout) * time.Second):
		klog.Error("Execution timeout!")
		w.WriteHeader(http.StatusGatewayTimeout)
		return
	}

	responseBytes, err := json.Marshal(response)
	if err != nil {
		klog.Error("500 internal error!")
	}

	requests.Inc()
	w.Header().Set("Content-Type", "application/json")
	w.Write(responseBytes)
}

func calculateFibonacci(number int, resChan chan int) {
	resChan <- fibonacci(number)
}

// Returns the fibonacci of a given number
func fibonacci(n int) int {
	if n < 2 {
		return n
	}

	return fibonacci(n-1) + fibonacci(n-2)
}
