package main

import (
	"os"
	"resource_manager/internal/cluster"
	"resource_manager/internal/consts"
	"resource_manager/internal/scaler"

	"k8s.io/klog"
)

func main() {
	klog.Info(consts.MSG_APP_STARTED)

	cluster.RegisterClientSet()

	scaler.NewScalerManager().Run()

	os.Exit(0)
}
