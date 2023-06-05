package main

import (
	"os"
	"resource_manager/internal/cluster"
	"resource_manager/internal/consts"

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

		// If node only has one pod with no utilization
		if podCount == 1 && node.GetCpuUtilization() == 0 {
			pod := node.ListPods()[0]
			if !pod.IsAlreadyWarm() {
				pod.WarmUp()
			}
		}

		if podCount == 1 && node.GetCpuUtilization() != 0 {
			pod := node.ListPods()[0]
			if pod.IsAlreadyWarm() {
				pod.UnsetWarm()
			}
		}
	}

	os.Exit(0)
}
