package prometheus

import (
	"os"

	"github.com/prometheus/client_golang/api"
	"k8s.io/klog"
)

var apiClient api.Client

// Initializes Prometheus client
func Init() {
	client, err := api.NewClient(api.Config{
		Address: os.Getenv("PROMETHEUS_URL"),
	})
	if err != nil {
		klog.Fatalf("Error creating client: %v\n", err)
	}

	klog.Info("Created Prometheus client.")

	apiClient = client
}
