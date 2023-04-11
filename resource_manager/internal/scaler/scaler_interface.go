package scaler

import "resource_manager/internal/cluster"

// Every scaler application should implement this interface
type ScalerInterface interface {
	getName() string
	shouldScale(cluster.ClusterMetrics) bool
	planScaling(cluster.ClusterMetrics)
	scale()
}
