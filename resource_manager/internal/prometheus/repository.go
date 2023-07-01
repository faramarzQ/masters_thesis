package prometheus

import (
	"context"
	"log"
	"time"

	v1 "github.com/prometheus/client_golang/api/prometheus/v1"
	"github.com/prometheus/common/model"
)

// Queries Prometheus for a given command and time duration
func Query(query string, from time.Time) model.Value {
	v1api := v1.NewAPI(apiClient)
	result, _, err := v1api.Query(context.Background(), query, from, v1.WithTimeout(5*time.Second))
	if err != nil {
		log.Fatalf("Error querying Prometheus: %v\n", err)
	}
	return result
}
