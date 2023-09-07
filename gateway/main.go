package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"k8s.io/klog/v2"
)

type response struct {
	Result      int    `json:"result"`
	ProcessedBy string `json:"processed_by"`
}

var (
	requests = prometheus.NewCounter(prometheus.CounterOpts{
		Name: "requests_total",
		Help: "Number of all requests.",
	})
)

func init() {
	prometheus.MustRegister(requests)
}

func main() {
	fmt.Println("Executing gateway.")

	http.HandleFunc("/", handler)
	http.Handle("/metrics", promhttp.Handler())

	err := http.ListenAndServe(":4444", nil)
	if err != nil {
		fmt.Println("Error occurred running server: ", err)
	}
}

// Handler for fibonacci calculator
func handler(w http.ResponseWriter, r *http.Request) {
	number, err := strconv.Atoi(r.URL.Query().Get("number"))
	klog.Info("Got a request: ", number)
	if err != nil {
		klog.Error("Wrong input!")
	}

	requests.Inc()

	ctx := context.Background()
	fibonacciHost := string(os.Getenv("FIBONACCI_NODEPORT_SERVICE_HOST"))
	fibonacciPort := os.Getenv("FIBONACCI_NODEPORT_SERVICE_PORT")
	url := fmt.Sprintf("http://%s:%s?number=%d", fibonacciHost, fibonacciPort, number)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		klog.Fatal("Error building http request with context: %s\n", err)
	}

	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		klog.Fatal("Error making http request: %s\n", err)
	}

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		klog.Fatal(err)
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(body)
}
