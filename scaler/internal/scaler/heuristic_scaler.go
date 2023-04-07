package scaler

import (
	"scaler/internal/cluster"
	"scaler/internal/consts"
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
	return false
}

func (hs *HeuristicScaler) scale(clusterMetrics cluster.ClusterMetrics) {

}
