package scaler

import (
	"math"
	"os"
	"resource_manager/internal/cluster"
	"resource_manager/internal/consts"
	"strconv"

	"k8s.io/klog"
)

type HeuristicScaler struct {
	baseScaler
}

func NewHeuristicScaler() *HeuristicScaler {
	return &HeuristicScaler{}
}

func (hs *HeuristicScaler) getName() string {
	return consts.HEURISTIC_SCALER
}

func (hs *HeuristicScaler) shouldScale(clusterMetrics cluster.ClusterMetrics) bool {
	klog.Info(consts.MSG_RUNNING_SHOULD_SCALE)
	defer klog.Info(consts.MSG_FINISHED_SHOULD_SCALE)

	// if int(clusterMetrics.GetAverageMemoryUtilization()) > consts.HEURISTIC_SCALER_UPPER_MEMORY_THRESHOLD {
	// 	return true
	// }

	offNodesCount := len(cluster.ListNodes().InClass(consts.OFF_CLASS))
	if offNodesCount == 0 {
		klog.Info("No Off class node found!")
		return false
	}

	upperCpuThreshold, err := strconv.Atoi(os.Getenv("HEURISTIC_SCALER_UPPER_CPU_THRESHOLD"))
	if err != nil {
		klog.Fatal(err)
	}

	averageCpuUtil := int(clusterMetrics.GetAverageCpuUtilization())

	klog.Info("ENV: Upper CPU threshold: ", upperCpuThreshold)
	klog.Info("Average CPU utilization: ", averageCpuUtil)

	if averageCpuUtil > upperCpuThreshold {
		return true
	}

	klog.Info("CPU utilization has not passed the threshold.")

	return false
}

func (hs *HeuristicScaler) planScaling(clusterMetrics cluster.ClusterMetrics) error {
	klog.Info(consts.MSG_RUNNING_SCALE_PLANNING)
	defer klog.Info(consts.MSG_FINISHED_SCALE_PLANNING)

	nodes := cluster.ListNodes()
	// hs.planScalingConsideringMemoryResource(clusterMetrics, nodes)
	err := hs.planScalingConsideringCpuResource(clusterMetrics, nodes)
	if err != nil {
		return err
	}

	return nil
}

// Plans scaling nodes considering their memory resource
func (hs *HeuristicScaler) planScalingConsideringMemoryResource(clusterMetrics cluster.ClusterMetrics, nodes cluster.NodeList) {

	// If utilization hasn't pass the threshold
	if int(clusterMetrics.GetAverageMemoryUtilization()) <= consts.HEURISTIC_SCALER_UPPER_MEMORY_THRESHOLD {
		return
	}

	// How much resource to add to satisfy the desired resource util
	usedMemory := (int(clusterMetrics.GetAverageMemoryUtilization()) * int(nodes.InClass(consts.ACTIVE_CLASS).TotalMemory())) / 100
	desiredMemory := (usedMemory / consts.HEURISTIC_SCALER_DESIRED_MEMORY_UTIL) * 100
	memoriesToAdd := desiredMemory - usedMemory

	nodesToTransit := hs.listNodesToSatisfyDesiredMemoryUtil(memoriesToAdd)

	if len(nodesToTransit) != 0 {
		var nodeTransition nodeTransition
		nodeTransition.from = consts.OFF_CLASS
		nodeTransition.to = consts.IDLE_CLASS
		nodeTransition.nodesList = nodesToTransit
		hs.setTransitions(nodeTransition)
	}
}

// Plans scaling nodes considering their CPU resource
func (hs *HeuristicScaler) planScalingConsideringCpuResource(clusterMetrics cluster.ClusterMetrics, nodes cluster.NodeList) error {

	// How much resource to add to satisfy the desired resource util
	desiredCpuThreshold, err := strconv.Atoi(os.Getenv("HEURISTIC_SCALER_DESIRED_CPU_UTIL"))
	if err != nil {
		klog.Error("Failed reading HEURISTIC_SCALER_DESIRED_CPU_UTIL from env")
		return err
	}

	usedCpu := (int(clusterMetrics.GetAverageCpuUtilization()) * int(nodes.InClass(consts.ACTIVE_CLASS).TotalCpu())) / 100
	desiredCpu := int(math.Ceil(float64(usedCpu * 100 / desiredCpuThreshold)))
	cpuCoresToAdd := desiredCpu - usedCpu

	klog.Info("CPU cores to add to meet the desired threshold: ", cpuCoresToAdd)

	nodesToTransit := hs.listNodesToSatisfyDesiredCpuUtil(cpuCoresToAdd)
	if len(nodesToTransit) != 0 {
		var nodeTransition nodeTransition
		nodeTransition.from = consts.OFF_CLASS
		nodeTransition.to = consts.ACTIVE_CLASS
		nodeTransition.nodesList = nodesToTransit
		hs.setTransitions(nodeTransition)
	}

	return nil
}

// Returns a list of nodes which their total memory is greater|equal to the given memory amount
func (hs *HeuristicScaler) listNodesToSatisfyDesiredMemoryUtil(memoriesToAdd int) cluster.NodeList {
	var desiredNodeList cluster.NodeList

	var tempMemory int
	for tempMemory < memoriesToAdd {
		efficientNode, ok := cluster.GetMostMemoryEfficientNode(desiredNodeList.Names(), consts.OFF_CLASS)
		if !ok {
			break
		}

		tempMemory += int(efficientNode.TotalMemory)
		desiredNodeList = append(desiredNodeList, *efficientNode)
	}

	return desiredNodeList
}

// Returns a list of nodes which their total CPU is greater|equal to the given CPU amount
func (hs *HeuristicScaler) listNodesToSatisfyDesiredCpuUtil(cpuCoresToAdd int) cluster.NodeList {
	var desiredNodeList cluster.NodeList

	var tempCpu int
	for tempCpu < cpuCoresToAdd {
		efficientNode, ok := cluster.GetMostCpuEfficientNode(desiredNodeList.Names(), consts.OFF_CLASS)
		if !ok {
			break
		}

		tempCpu += int(efficientNode.TotalCpu)
		desiredNodeList = append(desiredNodeList, *efficientNode)
	}

	return desiredNodeList
}
