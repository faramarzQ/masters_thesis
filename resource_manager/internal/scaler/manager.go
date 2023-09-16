package scaler

import (
	"resource_manager/internal/cluster"
	"resource_manager/internal/config"
	"resource_manager/internal/consts"

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
		NewFixedScaler(),
		NewHeuristicScaler(),
		NewProposedScaler(),
		NewSilencerScaler(),
	)

	return scalerManager
}

// Registers the active scaler into the scaler manager
func (sm *ScalerManager) RegisterActiveScaler(scalers ...ScalerInterface) {
	for _, scaler := range scalers {
		if scaler.getName() == config.ACTIVE_SCALER {
			sm.Scaler = scaler
			break
		}
	}
	klog.Infof(consts.MSG_REGISTERED_ACTIVE_SCALER, sm.Scaler.getName())
}

// Manages and executes the scheduler
func (sm ScalerManager) Run() {
	klog.Info(consts.MSG_RUNNING_SCALER_MANAGER)
	defer klog.Info(consts.MSG_FINISHED_SCALER_MANAGER)

	cluster.LabelNewNodes()

	clusterMetrics := cluster.GetClusterMetrics()

	sm.Scaler.prePlan()
	if sm.Scaler.shouldScale(clusterMetrics) {
		err := sm.Scaler.planScaling(clusterMetrics)
		if err != nil {
			sm.Scaler.onFail(err)
		}
		sm.Scaler.scale()
	}
}
