package cluster

import (
	"math"
	"os"
	"resource_manager/internal/consts"
	databaseModels "resource_manager/internal/database/model"
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
	ActiveClasses            []consts.NODE_CLASS
	NodesCount               int
	NodesDispersion          map[consts.NODE_CLASS]int
	PreviousState            string
	State                    string
	PreviousActionTaken      int8
	ActionTaken              int8
	EpsilonValue             uint8
	SuccessRequestRate       float64
	ClusterEnergyConsumption float64
	SuccessRateWeight        float32
	EnergyConsumptionWeight  float32
}

func GetClusterStatus() ClusterStatus {
	initialEpsilon, err := strconv.Atoi(os.Getenv("RL_INITIAL_EPSILON_VALUE"))
	if err != nil {
		klog.Fatal(err)
	}

	successRateWeight, err := strconv.Atoi(os.Getenv("RL_SUCCESS_RATE_WEIGHT"))
	if err != nil {
		klog.Fatal(err)
	}

	energyConsumptionWeight, err := strconv.Atoi(os.Getenv("RL_ENERGY_CONSUMPTION_WEIGHT"))
	if err != nil {
		klog.Fatal(err)
	}

	status := ClusterStatus{
		ActiveClasses:           consts.FUNCTIONING_CLASSES,
		NodesCount:              len(ListNodes()),
		NodesDispersion:         getsNodesDispersion(),
		EpsilonValue:            uint8(initialEpsilon),
		SuccessRateWeight:       float32(successRateWeight),
		EnergyConsumptionWeight: float32(energyConsumptionWeight),
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
	period, _ := strconv.Atoi(os.Getenv("PROMETHEUS_REQUESTS_PERIOD_MINUTE"))
	time := time.Now().Add(time.Duration(period) * time.Minute)
	result := prometheus.Query(consts.PROMETHEUS_METRIC_NAME_SUCCESS_REQUESTS, time)

	var successfulRequests int
	for _, vec := range result.(model.Vector) {
		successfulRequests += int(vec.Value)
	}
	return successfulRequests
}

// Returns tot
func getTotalRequests() int {
	period, _ := strconv.Atoi(os.Getenv("PROMETHEUS_REQUESTS_PERIOD_MINUTE"))
	time := time.Now().Add(time.Duration(period) * time.Minute)
	result := prometheus.Query(consts.PROMETHEUS_METRIC_NAME_TOTAL_REQUESTS, time)

	var totalRequests int
	for _, vec := range result.(model.Vector) {
		totalRequests += int(vec.Value)
	}
	return totalRequests
}

func GetSuccessRequestRate() float64 {
	totalRequests := getTotalRequests()
	if totalRequests == 0 {
		return 0
	}
	return float64(getSuccessfulRequests()) / float64(getTotalRequests())
}

// Calculates energy consumption of every node during the last scaling period
func CalculateEnergyConsumption(previousScalerExecutionLog databaseModels.ScalerExecutionLog) float64 {
	nodes := ListNodes().InClass(consts.ACTIVE_CLASS)
	from := previousScalerExecutionLog.CreatedAt
	minutesAgo := int(math.Floor(time.Now().Sub(previousScalerExecutionLog.CreatedAt).Seconds() / 30)) // every 30 second

	periodTimeSlots := []time.Time{}
	for i := 0; i <= minutesAgo; i++ {
		periodTimeSlots = append(periodTimeSlots, from.Add(time.Second*time.Duration(i*30)))
	}

	var energyConsumption float64
	var maxEnergyConsumption float64
	for _, node := range nodes {
		minPower := node.MinPowerConsumption
		maxPower := node.MaxPowerConsumption

		if minPower == 0 || maxPower == 0 {
			klog.Fatal("Minimum|Maximum power consumption is not set for node ", node.Name)
		}

		totalCpuCores := node.GetTotalCpuCores()

		var energyConsumptionOfNode float64
		var maxEnergyConsumptionOfNode float64
		for _, slot := range periodTimeSlots {
			usedCpuCoresInSlot, err := node.GetUsedCpuCoresAtGiveTime(slot)
			if err != nil {
				continue
			}

			cpuUtil := (usedCpuCoresInSlot / totalCpuCores) * 100
			powerAtSlot := (float64(maxPower-minPower) * cpuUtil / 100) + float64(minPower)
			energyAtSlot := powerAtSlot * (0.008333)
			maxEnergyAtSlot := float64(maxPower) * (0.008333)

			energyConsumptionOfNode += energyAtSlot
			maxEnergyConsumptionOfNode += maxEnergyAtSlot
		}

		energyConsumption += energyConsumptionOfNode
		maxEnergyConsumption += maxEnergyConsumptionOfNode
	}

	return energyConsumption / maxEnergyConsumption
}
