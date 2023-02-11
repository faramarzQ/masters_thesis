package app

import (
	"fmt"
	"scaler/internal/cluster"
	"scaler/internal/consts"
)

func Scale() {
	fmt.Println(consts.RUNNING_SCALER)

	cluster.RegisterNodes()

	// activeNodes := cluster.ListNodes()

	// fmt.Println(activeNodes)
}
