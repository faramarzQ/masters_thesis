package main

import (
	"os"
	"resource_manager/internal/cluster"
	"resource_manager/internal/consts"
	"resource_manager/internal/database"
	"resource_manager/internal/logger"
	"resource_manager/internal/prometheus"
	"resource_manager/internal/scaler"
	"time"

	"k8s.io/klog"
)

func main() {
	logger.Init()

	klog.Info(consts.MSG_SCALER_APP_STARTED)

	// initialize dependencies
	database.Init()
	prometheus.Init()
	cluster.RegisterClientSet()

	scaler.NewScalerManager().Run()

	klog.Info(consts.MSG_SCALER_APP_FINISHED)
	klog.Flush()

	time.Sleep(120 * time.Second)

	os.Exit(0)
}
