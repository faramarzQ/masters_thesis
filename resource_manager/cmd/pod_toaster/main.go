package main

import (
	"os"
	"resource_manager/internal/cluster"
	"resource_manager/internal/consts"
	"time"

	"k8s.io/klog"
)

func main() {
	klog.Info(consts.MSG_POD_TOASTER_APP_STARTED)

	cluster.RegisterClientSet()

	nodes := cluster.ListActiveNodes()

	for _, node := range nodes {
		podCount := len(node.ListPods())

		// If node has no pod
		if podCount == 0 {
			node.SetClass(consts.IDLE_CLASS)
		}

		if podCount > 1 {
			os.Exit(0)
		}

		// If node only has one pod with no utilization
		if node.GetCpuUtilization() == 0 {
			pod := node.ListPods()[0]

			if pod.IsAlreadyWarm() {
				// if has been warm for a while
				warmedAt := pod.GetWarmedAt()
				if warmedAt.Add(time.Minute*time.Duration(consts.WARM_POD_DURATION_MINUTES)).Unix() < time.Now().Unix() {
					node.SetClass(consts.IDLE_CLASS)
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

	os.Exit(0)
}
