package main

import (
	"os"
	"resource_manager/internal/cluster"
	"resource_manager/internal/consts"
	"resource_manager/internal/database"
	"resource_manager/internal/logger"
	"strconv"

	"k8s.io/klog"
)

var (
	nodeMaxPowerUpperThreshold, nodeMaxPowerLowerThreshold,
	nodeMinPowerUpperThreshold, nodeMinPowerLowerThreshold,
	desiredNumberOfGateways int
)

func main() {
	logger.Init()

	klog.Info(consts.MSG_TAGGER_APP_STARTED)

	database.Init()
	cluster.RegisterClientSet()

	readEnvVars()

	tagNodes()

	klog.Info(consts.MSG_TAGGER_APP_FINISHED)
	klog.Flush()

	os.Exit(0)
}

func readEnvVars() {
	var err error
	desiredNumberOfGateways, err = strconv.Atoi(os.Getenv("NUMBER_OF_GATEWAYS"))
	if err != nil {
		klog.Error("Failed reading env var: NUMBER_OF_GATEWAYS")
	}

	nodeMaxPowerUpperThreshold, err = strconv.Atoi(os.Getenv("NODE_MAX_POWER_UPPER_THRESHOLD"))
	if err != nil {
		klog.Error("Failed reading env var: NODE_MAX_POWER_UPPER_THRESHOLD")
	}

	nodeMaxPowerLowerThreshold, err = strconv.Atoi(os.Getenv("NODE_MAX_POWER_LOWER_THRESHOLD"))
	if err != nil {
		klog.Error("Failed reading env var: NODE_MAX_POWER_LOWER_THRESHOLD")
	}

	nodeMinPowerUpperThreshold, err = strconv.Atoi(os.Getenv("NODE_MIN_POWER_UPPER_THRESHOLD"))
	if err != nil {
		klog.Error("Failed reading env var: NODE_MIN_POWER_UPPER_THRESHOLD")
	}

	nodeMinPowerLowerThreshold, err = strconv.Atoi(os.Getenv("NODE_MIN_POWER_LOWER_THRESHOLD"))
	if err != nil {
		klog.Error("Failed reading env var: NODE_MIN_POWER_LOWER_THRESHOLD")
	}

}

func tagNodes() {
	nodes := cluster.ListNodes()

	if desiredNumberOfGateways <= len(nodes) {
		klog.Fatal("Number of gateway nodes is higher than number of nodes!")
	}

	var numberOfGateways int
	for _, node := range nodes {
		if node.IsMaster {
			continue
		}

		if numberOfGateways < desiredNumberOfGateways {
			node.SetGateway()
			node.UnsetWorker()
			numberOfGateways++
			continue
		}

		node.UnsetGateway()
		node.SetWorker()
	}
}
