package config

import (
	"os"
	"resource_manager/internal/consts"
	"strconv"
)

var (
	APP                 string
	ENV                 string
	CLUSTER_AUTH_CONFIG string
	CLUSTER_NAMESPACE   string
	ACTIVE_SCALER       string
	ACTIVE_SCHEDULER    string
)

func init() {
	initEnvironments()
}

func initEnvironments() {
	CLUSTER_NAMESPACE = os.Getenv("CLUSTER_NAMESPACE")

	ENV = os.Getenv("ENV")
	if os.Getenv("ENV") == consts.ENV_DEV_LOCAL {
		ENV = consts.ENV_DEV_LOCAL
		CLUSTER_AUTH_CONFIG = os.Getenv("CONFIG_DIR_DEV_LOCAL")
	} else {
		ENV = consts.ENV_DEV_MINIKUBE
		CLUSTER_AUTH_CONFIG = os.Getenv("CONFIG_DIR_DEV_MINIKUBE")
	}

	APP = os.Getenv("APP")
	appId, _ := strconv.Atoi(os.Getenv("APP_ID"))
	if APP == consts.SCALER_APPLICATION {
		ACTIVE_SCALER = consts.MAP_SCALER_ID_TO_NAME[uint(appId)]
	} else if APP == consts.SCHEDULER_APPLICATION {
		ACTIVE_SCHEDULER = consts.MAP_SCHEDULER_ID_TO_NAME[uint(appId)]
	}
}
