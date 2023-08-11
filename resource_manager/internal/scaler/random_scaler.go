package scaler

import (
	"math"
	"math/rand"
	"resource_manager/internal/cluster"
	"resource_manager/internal/consts"
	"strconv"
	"time"

	"k8s.io/klog"
)

type RandomScaler struct {
	baseScaler
}

func NewRandomScaler() *RandomScaler {
	return &RandomScaler{}
}

func (rs *RandomScaler) getName() string {
	return consts.RANDOM_SCALER
}

func (rs *RandomScaler) shouldScale(clusterMetrics cluster.ClusterMetrics) bool {
	klog.Info(consts.MSG_RUNNING_SHOULD_SCALE)
	defer klog.Info(consts.MSG_FINISHED_SHOULD_SCALE)

	rand.Seed(time.Now().UnixNano())
	randNum := rand.Intn(100)
	return randNum > 50
}

func (rs *RandomScaler) planScaling(clusterMetrics cluster.ClusterMetrics) error {
	klog.Info(consts.MSG_RUNNING_SCALE_PLANNING)
	defer klog.Info(consts.MSG_FINISHED_SCALE_PLANNING)

	nodes := cluster.ListNodes()

	var percentOfNodesToTransit float64 = 0.2
	var nodeTransition nodeTransition
	nodeTransition.from = consts.OFF_CLASS
	nodeTransition.to = consts.IDLE_CLASS

	// TODO: move to shouldScale
	offNodesCount := float64(len(nodes.InClass(consts.OFF_CLASS)))
	if offNodesCount == 0 {
		return nil
	}

	numberOfNodesToTransit := int64(math.Ceil(offNodesCount * percentOfNodesToTransit))
	nodeTransition.nodesList = cluster.GetRandomNodesFromNodeList(nodes.InClass(consts.OFF_CLASS), numberOfNodesToTransit)
	rs.setTransitions(nodeTransition)

	klog.Info("Number of nodes to scale: " + strconv.Itoa(len(nodeTransition.nodesList)))

	return nil
}
