package main

import (
	"os"
	"resource_manager/internal/cluster"
	"resource_manager/internal/consts"
	"resource_manager/internal/database"
	"resource_manager/internal/scaler"

	"k8s.io/klog"
)

func main() {
	klog.Info(consts.MSG_SCALER_APP_STARTED)

	database.Init()

	cluster.RegisterClientSet()

	scaler.NewScalerManager().Run()

	os.Exit(0)
}
