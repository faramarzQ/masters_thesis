package main

import (
	"os"
	"scaler/internal/cluster"
	"scaler/internal/consts"
	"scaler/internal/scaler"

	"k8s.io/klog"
)

func main() {
	klog.Info(consts.MSG_APP_STARTED)

	cluster.RegisterClientSet()

	scaler.NewScalerManager().Run()

	os.Exit(0)
}
