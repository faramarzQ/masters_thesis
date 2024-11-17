package scaler

import (
	"math"
	"math/rand"
	"os"
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

	return true

	rand.Seed(time.Now().UnixNano())
	randNum := rand.Intn(100)

	scalingProbability, err := strconv.Atoi(os.Getenv("RANDOM_SCALER_SCALING_PROBABILITY"))
	if err != nil {
		klog.Fatal(err)
	}

	klog.Info("Scaling probability: ", scalingProbability)
	klog.Info("Random number: ", randNum)

	return randNum < scalingProbability
}

func (rs *RandomScaler) planScaling(clusterMetrics cluster.ClusterMetrics) error {
	klog.Info(consts.MSG_RUNNING_SCALE_PLANNING)
	defer klog.Info(consts.MSG_FINISHED_SCALE_PLANNING)

	maxPercentOfNodesToTransit, err := strconv.ParseFloat(os.Getenv("RANDOM_SCALER_PERCENT_OF_NODES_TO_TRANSIT"), 64)
	if err != nil {
		return err
	}

	min := 9

	rand.Seed(time.Now().UnixNano())
	percentOfNodesToTransit := float64(rand.Intn(int(maxPercentOfNodesToTransit)-min)) + float64(min)

	klog.Info("Percent of nodes to transit: ", percentOfNodesToTransit)

	var nodeTransition nodeTransition
	nodeTransition.from = consts.OFF_CLASS
	nodeTransition.to = consts.ACTIVE_CLASS

	nodes := cluster.ListNodes()
	offNodesCount := float64(len(nodes.InClass(consts.OFF_CLASS)))
	if offNodesCount == 0 {
		klog.Info("No Off nodes found to be scaled")
		return nil
	}

	numberOfNodesToTransit := int64(math.Ceil(offNodesCount * percentOfNodesToTransit / 100))
	nodeTransition.nodesList = cluster.GetRandomNodesFromNodeList(nodes.InClass(consts.OFF_CLASS), numberOfNodesToTransit)
	rs.setTransitions(nodeTransition)

	klog.Info("Number of nodes to scale: " + strconv.Itoa(len(nodeTransition.nodesList)))

	return nil
}
