package consts

type SCALING_TYPE string
type NODE_CLASS string

var (
	SCALER_APPLICATION    = "scaler"
	SCHEDULER_APPLICATION = "scheduler"

	NODE_CLASS_LABEL_NAME     = "class"
	NODE_SCALED_AT_LABEL_NAME = "scaled_at"

	ACTIVE_CLASS NODE_CLASS = "active"
	IDLE_CLASS   NODE_CLASS = "idle"
	SLEEP_CLASS  NODE_CLASS = "sleep"
	OFF_CLASS    NODE_CLASS = "off"

	ENV_DEV_LOCAL    = "DEV_LOCAL"
	ENV_DEV_MINIKUBE = "DEV_MINIKUBE"
	ENV_PROD         = "PROD"

	MINIMUM_ACTIVE_NODES_COUNT = 2
	MINIMUM_IDLE_NODES_COUNT   = 0
	MINIMUM_SLEEP_NODES_COUNT  = 0

	FIXED_IDLE_NODES_COUNT = 1

	RANDOM_SCALER    = "Random scaler"
	HEURISTIC_SCALER = "Heuristic scaler"
	FIXED_SCALER     = "Fixed scaler"

	MAP_SCALER_ID_TO_NAME = map[uint]string{
		1: RANDOM_SCALER,
		2: FIXED_SCALER,
		3: HEURISTIC_SCALER,
	}

	HEURISTIC_SCALER_UPPER_MEMORY_THRESHOLD = 70
	HEURISTIC_SCALER_UPPER_CPU_THRESHOLD    = 70
	HEURISTIC_SCALER_DESIRED_MEMORY_UTIL    = 70
	HEURISTIC_SCALER_DESIRED_CPU_UTIL       = 70

	RANDOM_SCHEDULER = "RandomScheduler"

	MAP_SCHEDULER_ID_TO_NAME = map[uint]string{
		1: RANDOM_SCHEDULER,
	}

	SCALING_OUT SCALING_TYPE = "scaling out"
	SCALING_IN  SCALING_TYPE = "scaling in"

	RANDOM_SCHEDULER_NAME = "RandomScheduler"
)
