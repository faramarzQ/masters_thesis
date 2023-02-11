package config

import (
	"fmt"
	"os"
	"scaler/internal/consts"
)

var (
	ENV                 string
	CLUSTER_AUTH_CONFIG string
)

func init() {
	fmt.Println("here")
	if os.Getenv("ENV") == consts.ENV_DEV_LOCAL {
		ENV = consts.ENV_DEV_LOCAL
		CLUSTER_AUTH_CONFIG = os.Getenv("CONFIG_DIR_DEV_LOCAL")
	} else {
		ENV = consts.ENV_DEV_MINIKUBE
		CLUSTER_AUTH_CONFIG = os.Getenv("CONFIG_DIR_DEV_MINIKUBE")
	}
}
