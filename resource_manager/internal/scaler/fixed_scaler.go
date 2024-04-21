package scaler

import (
	"math"
	"os"
	"resource_manager/internal/cluster"
	"resource_manager/internal/consts"
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
	fixedNodesCount, _ := strconv.Atoi(os.Getenv("FIXED_IDLE_NODES_COUNT"))
	idleNodesCount := len(nodes.InClass(consts.IDLE_CLASS))

	klog.Info("Desired idle nodes: ", fixedNodesCount)
	klog.Info("Idle nodes: ", idleNodesCount)

	if idleNodesCount < fixedNodesCount {
		return true
	}

	offNodesCount := float64(len(nodes.InClass(consts.OFF_CLASS)))
	if offNodesCount == 0 {
		return false
	}

	return false
}

// In case any off node exists, transits them to the idle class to meet the fixed number required.
// Nodes are selected randomly
func (fs *FixedScaler) planScaling(cluster.ClusterMetrics) error {
	klog.Info(consts.MSG_RUNNING_SCALE_PLANNING)
	defer klog.Info(consts.MSG_FINISHED_SCALE_PLANNING)

	nodes := cluster.ListNodes()

	var nodeTransition nodeTransition
	nodeTransition.from = consts.OFF_CLASS
	nodeTransition.to = consts.IDLE_CLASS

	offNodesCount := float64(len(nodes.InClass(consts.OFF_CLASS)))

	fixedNodesCount, _ := strconv.Atoi(os.Getenv("FIXED_IDLE_NODES_COUNT"))
	idleNodesCount := len(nodes.InClass(consts.IDLE_CLASS))
	numberOfNodesToTransit := int64(math.Min(float64(fixedNodesCount-idleNodesCount), offNodesCount))
	nodeTransition.nodesList = cluster.GetRandomNodesFromNodeList(nodes.InClass(consts.OFF_CLASS), numberOfNodesToTransit)
	fs.setTransitions(nodeTransition)

	klog.Info("Number of nodes to scale: " + strconv.Itoa(len(nodeTransition.nodesList)))

	return nil
}
