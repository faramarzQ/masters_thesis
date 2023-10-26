package main

import (
	"resource_manager/internal/cluster"
	"resource_manager/internal/consts"
	"resource_manager/internal/prometheus"
	"resource_manager/internal/scheduler"

	"k8s.io/klog"
)

func main() {
	klog.Info(consts.MSG_SCHEDULER_APP_STARTED)

	cluster.RegisterClientSet()
	prometheus.Init()

	scheduler.NewSchedulerManager().Run()
}
