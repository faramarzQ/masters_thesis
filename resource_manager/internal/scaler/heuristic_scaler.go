package scaler

import (
	"resource_manager/internal/cluster"
	"resource_manager/internal/consts"

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

	if clusterMetrics.GetAverageMemoryUtilization() > 70 {
		return true
	}

	if clusterMetrics.GetAverageCpuUtilization() > 70 {
		return true
	}

	return false
}

func (hs *HeuristicScaler) planScaling(clusterMetrics cluster.ClusterMetrics) {
	klog.Info(consts.MSG_RUNNING_SHOULD_SCALE)
	defer klog.Info(consts.MSG_FINISHED_SHOULD_SCALE)

}

// func (hs *HeuristicScaler) scale() {

// }
