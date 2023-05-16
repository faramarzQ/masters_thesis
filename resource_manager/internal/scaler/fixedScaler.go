package scaler

import (
	"resource_manager/internal/cluster"
	"resource_manager/internal/consts"
	"resource_manager/internal/helpers"
	"strconv"

	"k8s.io/klog"
)

// The fixed scaler keeps the number of idle nods fixed (just by scaling out)
type FixedScaler struct {
	baseScaler
}

func NewFixedScaler() *FixedScaler {
	return &FixedScaler{}
}

func (fs *FixedScaler) getName() string {
	return consts.FIXED_SCALER
}

// Returns true if number of idle nodes is less than the desired fixed number
func (fs *FixedScaler) shouldScale(clusterMetrics cluster.ClusterMetrics) bool {
	klog.Info(consts.MSG_RUNNING_SHOULD_SCALE)
	defer klog.Info(consts.MSG_FINISHED_SHOULD_SCALE)

	nodes := cluster.ListNodes()

	idleNodesCount := len(nodes.InClass(consts.IDLE_CLASS))
	if idleNodesCount < consts.FIXED_IDLE_NODES_COUNT {
		return true
	}

	return false
}

// In case any off node exists, transits them the idle class to meet the fixed number required
func (fs *FixedScaler) planScaling(cluster.ClusterMetrics) {
	klog.Info(consts.MSG_RUNNING_SCALE_PLANNING)
	defer klog.Info(consts.MSG_FINISHED_SCALE_PLANNING)

	nodes := cluster.ListNodes()

	offNodesCount := float64(len(nodes.InClass(consts.OFF_CLASS)))
	if offNodesCount == 0 {
		return
	}

	var nodeTransition nodeTransition
	nodeTransition.from = consts.OFF_CLASS
	nodeTransition.to = consts.IDLE_CLASS

	idleNodesCount := len(nodes.InClass(consts.IDLE_CLASS))
	numberOfNodesToTransit := int64(consts.FIXED_IDLE_NODES_COUNT - idleNodesCount)
	nodeTransition.nodesList = helpers.GetRandomNodesFromNodeList(nodes.InClass(consts.OFF_CLASS), numberOfNodesToTransit)
	fs.setTransitions(nodeTransition)

	klog.Info("Number of nodes to scale: " + strconv.Itoa(len(nodeTransition.nodesList)))
}
