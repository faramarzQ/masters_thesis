package main

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"k8s.io/klog"
)

var (
	epsilonValue = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "epsilon_value",
	})

	energyConsumption = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "energy_consumption",
	})

	reward = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "reward",
	})
)

func init() {
	prometheus.MustRegister(epsilonValue)
	prometheus.MustRegister(energyConsumption)
	prometheus.MustRegister(reward)
}

func main() {
	fmt.Println("Executing monitoring server.")

	http.Handle("/metrics", promhttp.Handler())
	http.HandleFunc("/epsilon", epsilonHandler)
	http.HandleFunc("/energy_consumption", energyConsumptionHandler)
	http.HandleFunc("/reward", rewardHandler)

	err := http.ListenAndServe(":5555", nil)
	if err != nil {
		fmt.Println("Error occurred running server: ", err)
	}
}

func epsilonHandler(w http.ResponseWriter, r *http.Request) {
	number, err := strconv.ParseFloat(r.URL.Query().Get("value"), 64)
	klog.Info("Received epsilon value: ", number)
	if err != nil {
		klog.Error("Wrong input!")
	}

	epsilonValue.Set(number)
}

func energyConsumptionHandler(w http.ResponseWriter, r *http.Request) {
	number, err := strconv.ParseFloat(r.URL.Query().Get("value"), 64)
	klog.Info("Received energy consumption value: ", number)
	if err != nil {
		klog.Error("Wrong input!")
	}

	energyConsumption.Set(number)
}

func rewardHandler(w http.ResponseWriter, r *http.Request) {
	number, err := strconv.ParseFloat(r.URL.Query().Get("value"), 64)
	klog.Info("Received a reward value: ", number)
	if err != nil {
		klog.Error("Wrong input!")
	}

	reward.Set(number)
}
