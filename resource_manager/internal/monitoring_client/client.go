package monitoring_client

import (
	"fmt"
	"net/http"
	"os"

	"k8s.io/klog"
)

func SetEpsilonValue(epsilon float64) {
	_, err := http.Get(os.Getenv("MONITORING_SERVER_URL") + "/epsilon?value=" + fmt.Sprintf("%f", epsilon))
	if err != nil {
		klog.Error("Failed setting epsilon value to monitoring server.")
	}
}

func SetEnergyConsumptionValue(energyConsumption float64) {
	_, err := http.Get(os.Getenv("MONITORING_SERVER_URL") + "/energy_consumption?value=" + fmt.Sprintf("%f", energyConsumption))
	if err != nil {
		klog.Error("Failed setting energy consumption value to monitoring server.")
	}
}

func SetRewardValue(reward float64) {
	_, err := http.Get(os.Getenv("MONITORING_SERVER_URL") + "/reward?value=" + fmt.Sprintf("%f", reward))
	if err != nil {
		klog.Error("Failed setting reward value to monitoring server.")
	}
}
