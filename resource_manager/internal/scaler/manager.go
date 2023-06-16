package scaler

import (
	"resource_manager/internal/cluster"
	"resource_manager/internal/config"
	"resource_manager/internal/consts"
	"resource_manager/internal/database/repository"

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
		NewSilencerScaler(),
		NewProposedScaler(),
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

	repository.InsertScalerExecutionLog(sm.Scaler.getName())

	cluster.LabelNewNodes()

	clusterMetrics := cluster.GetClusterMetrics()

	if sm.Scaler.shouldScale(clusterMetrics) {
		sm.Scaler.planScaling(clusterMetrics)
		sm.Scaler.scale()
		sm.postScale()
	}
}

func (sm ScalerManager) postScale() {
	// store executed scaler's name in the cluster
	cluster.MasterNode().SetAnnotation(consts.ACTIVE_SCALER_LABEL_NAME, sm.Scaler.getName())
}
