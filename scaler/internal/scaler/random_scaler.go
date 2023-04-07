package scaler

import (
	"fmt"
	"math"
	"math/rand"
	"scaler/internal/cluster"
	"scaler/internal/consts"
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
	return randNum > 30
}

func (rs *RandomScaler) planScaling(clusterMetrics cluster.ClusterMetrics) {
	klog.Info(consts.MSG_RUNNING_SHOULD_SCALE)
	defer klog.Info(consts.MSG_FINISHED_SHOULD_SCALE)

	scalingType := rs.calculateScalingType()
	klog.Info("Scaling type: \"" + scalingType + "\"")

	nodes := cluster.ListNodes()

	var percentOfNodesToTransit float64 = 0.2
	var nodeTransition nodeTransition
	if scalingType == consts.SCALING_IN {
		nodeTransition.from = consts.ACTIVE_CLASS
		nodeTransition.to = consts.OFF_CLASS

		activeNodesCount := float64(len(nodes.InClass(consts.ACTIVE_CLASS)))
		if activeNodesCount == 0 {
			return
		}
		numberOfNodesToTransit := int64(math.Ceil(activeNodesCount * percentOfNodesToTransit))
		nodeTransition.nodesList = rs.getRandomNodes(nodes.InClass(consts.ACTIVE_CLASS), numberOfNodesToTransit)
	} else {
		nodeTransition.from = consts.OFF_CLASS
		nodeTransition.to = consts.ACTIVE_CLASS

		offNodesCount := float64(len(nodes.InClass(consts.OFF_CLASS)))
		if offNodesCount == 0 {
			return
		}
		numberOfNodesToTransit := int64(math.Ceil(offNodesCount * percentOfNodesToTransit))
		nodeTransition.nodesList = rs.getRandomNodes(nodes.InClass(consts.OFF_CLASS), numberOfNodesToTransit)
	}
	klog.Info("Number of nodes to scale: " + strconv.Itoa(len(nodeTransition.nodesList)))

	rs.setTransitions(nodeTransition)
}

func (rs *RandomScaler) getRandomNodes(nodeList cluster.NodeList, numberOfNodesToSelect int64) cluster.NodeList {
	nodesIndexesToSelect := []int{}

	for len(nodesIndexesToSelect) < int(numberOfNodesToSelect) {
		rand.Seed(time.Now().UnixNano())
		randomNum := rand.Intn(int(len(nodeList)))

		var nodeAlreadySelected bool
		for i := 0; i < len(nodesIndexesToSelect); i++ {
			if randomNum == nodesIndexesToSelect[i] {
				nodeAlreadySelected = true
				break
			}
		}

		fmt.Println(randomNum)
		if !nodeAlreadySelected {
			nodesIndexesToSelect = append(nodesIndexesToSelect, randomNum)
		}
	}

	var nodesToSelect cluster.NodeList
	for i := 0; i < len(nodesIndexesToSelect); i++ {
		nodesToSelect = append(nodesToSelect, nodeList[nodesIndexesToSelect[i]])
	}

	return nodesToSelect
}

func (rs *RandomScaler) calculateScalingType() consts.SCALING_TYPE {
	rand.Seed(time.Now().UnixNano())
	randNum := rand.Intn(100)

	if randNum > 50 {
		return consts.SCALING_IN
	}

	return consts.SCALING_OUT
}
