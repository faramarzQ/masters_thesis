package prometheus

import (
	"log"
	"os"

	"github.com/prometheus/client_golang/api"
)

var apiClient api.Client

// Initializes Prometheus client
func Init() {
	client, err := api.NewClient(api.Config{
		Address: os.Getenv("PROMETHEUS_URL"),
	})
	if err != nil {
		log.Fatalf("Error creating client: %v\n", err)
	}

	apiClient = client
}
