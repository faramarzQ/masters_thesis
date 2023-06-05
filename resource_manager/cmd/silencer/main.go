package main

import (
	"fmt"
	"os"
	"resource_manager/internal/cluster"
	"resource_manager/internal/consts"
	"time"

	"k8s.io/klog"
)

func main() {
	klog.Info(consts.MSG_POD_TOASTER_APP_STARTED)

	cluster.RegisterClientSet()

	// silenceActiveNodes()
	silenceIdleNodes()

	os.Exit(0)
}

// For every active node, checks if it has been inactive for a while then silences them to lower level classes
func silenceActiveNodes() {
	nodes := cluster.ListActiveNodes()

	// for every scaler type, directly change nodes to off class,
	// but for the proposed scaler, the idle class is also usable
	activeScaler := cluster.MasterNode().Labels[consts.ACTIVE_SCALER_LABEL_NAME]
	targetClass := consts.IDLE_CLASS
	if activeScaler == "proposed" {
		targetClass = consts.OFF_CLASS
	}

	fmt.Println("here")
	for _, node := range nodes {
		fmt.Println("here")
		fmt.Println(node.Annotations)
		// node.SetClass(targetClass)
		return
		podCount := len(node.ListPods())

		if podCount > 1 {
			os.Exit(0)
		}

		// If node has no pod
		if podCount == 0 {
			node.SetClass(consts.IDLE_CLASS)
		}

		// If node only has one pod with no utilization
		if node.GetCpuUtilization() == 0 {
			pod := node.ListPods()[0]

			if pod.IsAlreadyWarm() {
				// if has been warm for a while
				warmedAt := pod.GetWarmedAt()
				if warmedAt.Add(time.Minute*time.Duration(consts.WARM_POD_DURATION_MINUTES)).Unix() < time.Now().Unix() {
					node.SetClass(targetClass)
					pod.UnsetWarm()
				}

				return
			} else {
				pod.WarmUp()
			}
		}

		if node.GetCpuUtilization() != 0 {
			pod := node.ListPods()[0]
			if pod.IsAlreadyWarm() {
				pod.UnsetWarm()
			}
		}
	}
}

// For every idle node, checks if it has been idle for a when, then silences them to off class
func silenceIdleNodes() {
	// Only available for the proposed scaler
	activeScaler := cluster.MasterNode().Labels[consts.ACTIVE_SCALER_LABEL_NAME]
	if activeScaler != "proposed" {
		return
	}

	nodes := cluster.ListNodes().InClass(consts.IDLE_CLASS)

	for _, node := range nodes {
		// if node has been idle for a specific time, scale to off
		scaledAt := node.GetScaledAt()
		if scaledAt.Add(time.Minute*time.Duration(consts.IDLE_NODE_DURATION_MINUTES)).Unix() < time.Now().Unix() {
			node.SetClass(consts.OFF_CLASS)
		}
	}
}
