package scaler

import (
	"scaler/internal/cluster"
	"scaler/internal/consts"

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

func (bs *baseScaler) setTransitions(transitions ...nodeTransition) {
	bs.nodeTransitions = transitions
}

func (bs *baseScaler) scale() {
	klog.Info(consts.MSG_RUNNING_SCALER)
	defer klog.Info(consts.MSG_FINISHED_SCALER)

	for i := 0; i < len(bs.nodeTransitions); i++ {
		toClass := bs.nodeTransitions[i].to
		for j := 0; j < len(bs.nodeTransitions[i].nodesList); j++ {
			node := bs.nodeTransitions[i].nodesList[j]
			node.SetClass(toClass)
			klog.Info("Transitioned \"" + node.Name + "\" From \"" + string(bs.nodeTransitions[i].from) + "\" To \"" + string(toClass) + "\"")
		}
	}
}
