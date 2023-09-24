package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"k8s.io/klog/v2"
)

type response struct {
	Result      int    `json:"result"`
	ProcessedBy string `json:"processed_by"`
}

var (
	buckets  = []float64{.025, .05, .1, .25, .5, 1, 2.5, 5, 10, 30, 60}
	requests = prometheus.NewCounter(prometheus.CounterOpts{
		Name: "requests_total",
		Help: "Number of all requests.",
	})
	requestTimeouts = prometheus.NewCounter(prometheus.CounterOpts{
		Name: "request_timeout_total",
		Help: "Number of all timeout requests.",
	})
	responseTime = prometheus.NewHistogramVec(prometheus.HistogramOpts{
		Name:    "response_duration",
		Help:    "Time duration of responses.",
		Buckets: buckets,
	}, []string{})

	fibonacciHost, fibonacciPort string
)

func init() {
	prometheus.MustRegister(requests)
	prometheus.MustRegister(responseTime)
}

func main() {
	fmt.Println("Executing gateway.")

	http.HandleFunc("/", handler)
	http.Handle("/metrics", promhttp.Handler())

	fibonacciHost = string(os.Getenv("FIBONACCI_NODEPORT_SERVICE_HOST"))
	fibonacciPort = os.Getenv("FIBONACCI_NODEPORT_SERVICE_PORT")

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
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	requests.Inc()

	start := time.Now()

	ctx := context.Background()
	url := fmt.Sprintf("http://%s:%s?number=%d", fibonacciHost, fibonacciPort, number)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		klog.Error("Error building http request with context: %s\n", err)
	}

	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		klog.Error("Error making http request: ", err.Error())
		return
	}

	duration := time.Since(start)
	responseTime.WithLabelValues().Observe(duration.Seconds())

	if res.StatusCode == http.StatusGatewayTimeout {
		requestTimeouts.Inc()
		w.WriteHeader(http.StatusGatewayTimeout)
		return
	}

	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		klog.Error(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(body)
	return
}
