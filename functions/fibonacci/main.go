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
	fmt.Println("Listening on port 3333")
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

	result := <-resChan
	if result == -1 {
		w.WriteHeader(http.StatusGatewayTimeout)
		return
	}

	res := response{
		ProcessedBy: os.Getenv("HOSTNAME"),
		Result:      result,
	}

	responseBytes, err := json.Marshal(res)
	if err != nil {
		klog.Error("500 internal error!")
	}

	requests.Inc()

	w.Header().Set("Content-Type", "application/json")
	w.Write(responseBytes)
	return
}

func calculateFibonacci(number int, resChan chan int) {
	timeoutTime := time.Now().Add(time.Second * time.Duration(timeout))
	resChan <- fibonacci(number, timeoutTime, resChan)
}

// Returns the fibonacci of a given number
func fibonacci(n int, timeout time.Time, resChan chan int) int {
	if time.Now().After(timeout) {
		resChan <- -1
		return -1
	}

	if n < 2 {
		return n
	}

	return fibonacci(n-1, timeout, resChan) + fibonacci(n-2, timeout, resChan)
}
