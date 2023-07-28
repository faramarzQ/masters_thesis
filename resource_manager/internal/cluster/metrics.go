package cluster

import (
	"os"
	"resource_manager/internal/consts"
	"resource_manager/internal/prometheus"
	"strconv"
	"time"

	"github.com/prometheus/common/model"
	"k8s.io/klog/v2"
)

type ClusterMetrics struct {
	ActiveNodesMetrics ActiveNodesMetrics
}

// Metrics of cluster nodes
type NodeMetrics struct {
	CpuUtilization    float64
	MemoryUtilization float64
}

// Calculated metric values for on every active nodes
type ActiveNodesMetrics map[string]NodeMetrics

type ClusterStatus struct {
	ActiveClasses        []consts.NODE_CLASS
	NodesCount           int
	NodesDispersion      map[consts.NODE_CLASS]int
	SuccessfulRequests   int
	TotalRequests        int
	PreviousState        string
	PreviousActionTaken  int8
	PreviousEpsilonValue uint8
	ExecutedPreviously   bool
}

func GetClusterStatus() ClusterStatus {
	initialEpsilon, err := strconv.Atoi(os.Getenv("RL_INITIAL_EPSILON_VALUE"))
	if err != nil {
		klog.Fatal(err)
	}

	status := ClusterStatus{
		ActiveClasses:        consts.FUNCTIONING_CLASSES,
		NodesCount:           len(ListNodes()),
		NodesDispersion:      getsNodesDispersion(),
		SuccessfulRequests:   getSuccessfulRequests(),
		TotalRequests:        getTotalRequests(),
		PreviousEpsilonValue: uint8(initialEpsilon),
	}

	return status
}

func GetClusterMetrics() ClusterMetrics {
	return ClusterMetrics{
		GetActiveNodesMetrics(),
	}
}

func GetActiveNodesMetrics() ActiveNodesMetrics {
	activeNodes := ListActiveNodes()
	var activeNodesMetrics ActiveNodesMetrics = make(ActiveNodesMetrics, len(activeNodes))
	for _, node := range activeNodes {
		cpuUtilization := node.GetCpuUtilization()
		memoryUtilization := node.GetMemoryUtilization()

		activeNodesMetrics[node.Name] = NodeMetrics{
			CpuUtilization:    cpuUtilization,
			MemoryUtilization: memoryUtilization,
		}
	}
	return activeNodesMetrics
}

func (cm ClusterMetrics) GetAverageCpuUtilization() float64 {
	var sumCpuUtil float64
	for _, resourceMetrics := range cm.ActiveNodesMetrics {
		sumCpuUtil += resourceMetrics.CpuUtilization
	}
	return sumCpuUtil / float64(len(cm.ActiveNodesMetrics))
}

func (cm ClusterMetrics) GetAverageMemoryUtilization() float64 {
	var sumMemoryUtil float64
	for _, resourceMetrics := range cm.ActiveNodesMetrics {
		sumMemoryUtil += resourceMetrics.MemoryUtilization
	}
	return sumMemoryUtil / float64(len(cm.ActiveNodesMetrics))
}

// Returns number of successful requests
func getSuccessfulRequests() int {
	period, _ := strconv.Atoi(os.Getenv("PROMETHEUS_SUCCESS_REQUESTS_PERIOD_MINUTE"))
	time := time.Now().Add(time.Duration(period) * time.Minute)
	result := prometheus.Query(os.Getenv("PROMETHEUS_METRIC_NAME_SUCCESS_REQUESTS"), time)

	var successfulRequests int
	for _, vec := range result.(model.Vector) {
		successfulRequests += int(vec.Value)
	}
	return successfulRequests
}

// Returns tot
func getTotalRequests() int {
	period, _ := strconv.Atoi(os.Getenv("PROMETHEUS_TOTAL_REQUESTS_PERIOD_MINUTE"))
	time := time.Now().Add(time.Duration(period) * time.Minute)
	result := prometheus.Query(os.Getenv("PROMETHEUS_METRIC_NAME_TOTAL_REQUESTS"), time)

	var totalRequests int
	for _, vec := range result.(model.Vector) {
		totalRequests += int(vec.Value)
	}
	return totalRequests
}
