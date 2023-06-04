package scaler

import (
	"math"
	"resource_manager/internal/cluster"
	"resource_manager/internal/consts"
	"resource_manager/internal/helpers"
	"strconv"

	"k8s.io/klog"
)

type SilencerScaler struct {
	baseScaler
}

func NewSilencerScaler() *SilencerScaler {
	return &SilencerScaler{}
}

func (rs *SilencerScaler) getName() string {
	return consts.SILENCER_SCALER
}

func (rs *SilencerScaler) planScaling(clusterMetrics cluster.ClusterMetrics) {
	klog.Info(consts.MSG_RUNNING_SCALE_PLANNING)
	defer klog.Info(consts.MSG_FINISHED_SCALE_PLANNING)

	nodes := cluster.ListNodes()

	// for all nodes:
	// 	if has one pod
	// 		if util is zer
	// 			if has idle_at label and is greater than CONST value:
	// 				turn off

	// 		else
	// 			remove idle_at label

	// 	else
	// 		remove idle_at label
	var percentOfNodesToTransit float64 = 0.2
	var nodeTransition nodeTransition
	nodeTransition.from = consts.OFF_CLASS
	nodeTransition.to = consts.IDLE_CLASS

	offNodesCount := float64(len(nodes.InClass(consts.OFF_CLASS)))
	if offNodesCount == 0 {
		return
	}

	numberOfNodesToTransit := int64(math.Ceil(offNodesCount * percentOfNodesToTransit))
	nodeTransition.nodesList = helpers.GetRandomNodesFromNodeList(nodes.InClass(consts.OFF_CLASS), numberOfNodesToTransit)
	rs.setTransitions(nodeTransition)

	klog.Info("Number of nodes to scale: " + strconv.Itoa(len(nodeTransition.nodesList)))
}
