package scaler

import (
	"resource_manager/internal/cluster"
	"resource_manager/internal/config"
	"resource_manager/internal/consts"
	"resource_manager/internal/database/model"
	"resource_manager/internal/database/repository"

	"k8s.io/klog"
)

type baseScaler struct {
	nodeTransitions []nodeTransition
}

type nodeTransition struct {
	from      consts.NODE_CLASS
	to        consts.NODE_CLASS
	nodesList cluster.NodeList
}

var (
	scalerExecutionLog         model.ScalerExecutionLog
	previousScalerExecutionLog *model.ScalerExecutionLog
)

func (bs *baseScaler) setTransitions(transitions ...nodeTransition) {
	bs.nodeTransitions = append(bs.nodeTransitions, transitions...)
}

func (bs *baseScaler) shouldScale(_ cluster.ClusterMetrics) bool {
	return true
}

func (bs *baseScaler) onFail(err error) {
	klog.Info(consts.MSG_FAILED_PLANNING)
	repository.DeleteScalingExecutionLog(&scalerExecutionLog)
	klog.Fatal(err)
}

func (bs *baseScaler) prePlan() {
	klog.Info(consts.MSG_RUNNING_PRE_PLANNING)
	defer klog.Info(consts.MSG_FINISHED_PRE_PLANNING)

	cluster.MasterNode().SetAnnotation(consts.ACTIVE_SCALER_LABEL_NAME, config.ACTIVE_SCALER)

	// Fetch previous execution log
	previousScalerExecutionLog = repository.GetPreviousScalerExecutionLog(config.ACTIVE_SCALER)

	// Log execution
	var err error
	scalerExecutionLog, err = repository.InsertScalerExecutionLog(previousScalerExecutionLog, config.ACTIVE_SCALER)
	if err != nil {
		klog.Fatal(err)
	}
}

func (bs *baseScaler) scale() {
	klog.Info(consts.MSG_RUNNING_SCALER)
	defer klog.Info(consts.MSG_FINISHED_SCALER)

	for i := 0; i < len(bs.nodeTransitions); i++ {
		toClass := bs.nodeTransitions[i].to
		for j := 0; j < len(bs.nodeTransitions[i].nodesList); j++ {
			node := bs.nodeTransitions[i].nodesList[j]
			node.SetClass(toClass)
			// repository.InsertScalingLog(scalerExecutionLog, node.Name, toClass)

			klog.Info("Transitioned \"" + node.Name + "\" From \"" + string(bs.nodeTransitions[i].from) + "\" To \"" + string(toClass) + "\"")
		}

		// Only active nodes should have a function
		if toClass != consts.ACTIVE_CLASS {
			for j := 0; j < len(bs.nodeTransitions[i].nodesList); j++ {
				node := bs.nodeTransitions[i].nodesList[j]
				node.RemovePods()
			}
		}
	}
}
