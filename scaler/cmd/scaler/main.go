package main

import (
	"fmt"
	"scaler/internal/app"
	"scaler/internal/cluster"
	"scaler/internal/consts"
)

func main() {
	fmt.Println(consts.APP_STARTED)

	cluster.RegisterClientSet()

	app.Scale()
}
