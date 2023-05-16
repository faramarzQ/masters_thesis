package scaler

import (
	"math"
	"resource_manager/internal/cluster"
	"resource_manager/internal/consts"
	"resource_manager/internal/helpers"

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

	if int(clusterMetrics.GetAverageMemoryUtilization()) > consts.HEURISTIC_SCALER_UPPER_MEMORY_THRESHOLD {
		return true
	}

	if int(clusterMetrics.GetAverageCpuUtilization()) > consts.HEURISTIC_SCALER_UPPER_CPU_THRESHOLD {
		return true
	}

	return false
}

func (hs *HeuristicScaler) planScaling(clusterMetrics cluster.ClusterMetrics) {
	klog.Info(consts.MSG_RUNNING_SCALE_PLANNING)
	defer klog.Info(consts.MSG_FINISHED_SCALE_PLANNING)

	nodes := cluster.ListNodes()

	offNodesCount := len(nodes.InClass(consts.OFF_CLASS))
	if offNodesCount == 0 {
		return
	}

	if int(clusterMetrics.GetAverageMemoryUtilization()) > consts.HEURISTIC_SCALER_UPPER_MEMORY_THRESHOLD {
		var nodeTransition nodeTransition
		nodeTransition.from = consts.OFF_CLASS
		nodeTransition.to = consts.IDLE_CLASS

		var percentOfNodesToTransit float64 = 0.2
		numberOfNodesToTransit := int64(math.Ceil(float64(len(nodes.InClass(consts.ACTIVE_CLASS))) * percentOfNodesToTransit))
		nodeTransition.nodesList = helpers.GetRandomNodesFromNodeList(nodes.InClass(consts.OFF_CLASS), numberOfNodesToTransit)
		hs.setTransitions(nodeTransition)
	}

	if int(clusterMetrics.GetAverageCpuUtilization()) > consts.HEURISTIC_SCALER_UPPER_CPU_THRESHOLD {
		var nodeTransition nodeTransition
		nodeTransition.from = consts.OFF_CLASS
		nodeTransition.to = consts.IDLE_CLASS

		var percentOfNodesToTransit float64 = 0.2
		numberOfNodesToTransit := int64(math.Ceil(float64(len(nodes.InClass(consts.ACTIVE_CLASS))) * percentOfNodesToTransit))
		nodeTransition.nodesList = helpers.GetRandomNodesFromNodeList(nodes.InClass(consts.OFF_CLASS), numberOfNodesToTransit)
		hs.setTransitions(nodeTransition)
	}
}
