package config

import (
	"flag"
	"os"
	"scaler/internal/consts"
)

var (
	ENV                 string
	CLUSTER_AUTH_CONFIG string
	CLUSTER_NAMESPACE   string
	ACTIVE_SCALER       string
)

func init() {
	initEnvironments()
	initFlags()
}

func initEnvironments() {
	if os.Getenv("ENV") == consts.ENV_DEV_LOCAL {
		ENV = consts.ENV_DEV_LOCAL
		CLUSTER_AUTH_CONFIG = os.Getenv("CONFIG_DIR_DEV_LOCAL")
	} else {
		ENV = consts.ENV_DEV_MINIKUBE
		CLUSTER_AUTH_CONFIG = os.Getenv("CONFIG_DIR_DEV_MINIKUBE")
	}

	CLUSTER_NAMESPACE = os.Getenv("CLUSTER_NAMESPACE")
}

func initFlags() {
	var scalerId uint
	flag.UintVar(&scalerId, "scaler", 1, "Set to id number of the scaler application to run. Random scaler is the default")

	flag.Parse()
	ACTIVE_SCALER = consts.MAP_SCALER_ID_TO_NAME[scalerId]
}
