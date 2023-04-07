package scaler

import (
	"scaler/config"
	"scaler/internal/cluster"
	"scaler/internal/consts"

	"k8s.io/klog"
)

type ScalerManager struct {
	Scaler ScalerInterface
}

func NewScalerManager() *ScalerManager {
	scalerManager := &ScalerManager{}

	// Register your new scaler here
	scalerManager.RegisterActiveScaler(
		NewRandomScaler(),
		// NewHeuristicScaler(),
	)

	return scalerManager
}

// Registers the active scaler into the scaler manager
func (sm *ScalerManager) RegisterActiveScaler(scalers ...ScalerInterface) {
	for _, scaler := range scalers {
		if scaler.getName() == config.ACTIVE_SCALER {
			sm.Scaler = scaler
		}
	}
	klog.Infof(consts.MSG_REGISTERED_ACTIVE_SCALER, sm.Scaler.getName())
}

func (sm ScalerManager) Run() {
	klog.Info(consts.MSG_RUNNING_SCALER_MANAGER)
	defer klog.Info(consts.MSG_FINISHED_SCALER_MANAGER)

	cluster.LabelNewNodes()

	clusterMetrics := cluster.GetClusterMetrics()

	if sm.Scaler.shouldScale(clusterMetrics) {
		sm.Scaler.planScaling(clusterMetrics)
	}

	sm.Scaler.scale()
}
