PROJECT_ROOT := $(realpath $(dir $(firstword $(MAKEFILE_LIST))))
SILENCER_ENTRY := $(PROJECT_ROOT)/cmd/silencer/main.go
SCALER_ENTRY := $(PROJECT_ROOT)/cmd/scaler/main.go
SCHEDULER_ENTRY := $(PROJECT_ROOT)/cmd/scheduler/main.go

BIN_DIR := $(PROJECT_ROOT)/bin
SCHEDULER_OUTPUT := $(BIN_DIR)/scheduler
SCALER_OUTPUT := $(BIN_DIR)/scaler
SILENCER_OUTPUT := $(BIN_DIR)/silencer

SCRIPTS_DIR := $(PROJECT_ROOT)/scripts
ENV_LOCAL := env_local.sh
ENV_MINIKUBE := env_minikube.sh
ENV_MINIKUBE_BIN := env_minikube_bin.sh

SCHEDULER_DEFAULT_FLAGS=--authentication-kubeconfig=/etc/kubernetes/scheduler.conf --authorization-kubeconfig=/etc/kubernetes/scheduler.conf --config=../config/scheduler/random-scheduler-config.yml --secure-port=10260

# Scaler

run-random-scaler:
	cd $(SCRIPTS_DIR) && . ./$(ENV_LOCAL) \
	&& APP=scaler APP_ID=1 go run $(SCALER_ENTRY) $(flags)

build-scaler:
	go build -o $(SCALER_OUTPUT) $(SCALER_ENTRY) $(flags)

run-random-scaler-bin:
	cd $(SCRIPTS_DIR) && . ./$(ENV_MINIKUBE_BIN) \
	&& cd $(BIN_DIR) && APP=scaler APP_ID=1 ./scaler $(flags)

run-fixed-scaler-bin:
	cd $(SCRIPTS_DIR) && . ./$(ENV_MINIKUBE_BIN) \
	&& cd $(BIN_DIR) && APP=scaler APP_ID=2 ./scaler $(flags)

run-heuristic-scaler-bin:
	cd $(SCRIPTS_DIR) && . ./$(ENV_MINIKUBE_BIN) \
	&& cd $(BIN_DIR) && APP=scaler APP_ID=3 ./scaler $(flags)


# Scheduler

run-random-scheduler:
	cd $(SCRIPTS_DIR) && . ./$(ENV_LOCAL) \
	&& APP=scheduler APP_ID=1 go run $(SCHEDULER_ENTRY) $(flags)

build-scheduler:
	go build -o $(SCHEDULER_OUTPUT) $(SCHEDULER_ENTRY) $(flags)

run-random-scheduler-bin:
	cd $(SCRIPTS_DIR) && . ./$(ENV_MINIKUBE_BIN) \
	&& cd $(BIN_DIR) && APP=scheduler APP_ID=1 ./scheduler $(SCHEDULER_DEFAULT_FLAGS) $(flags)

# Silencer

build-silencer:
		go build -o $(SILENCER_OUTPUT) $(SILENCER_ENTRY) $(flags)

run-silencer:
	cd $(SCRIPTS_DIR) && . ./$(ENV_MINIKUBE_BIN) \
	&& cd $(BIN_DIR) && ./silencer $(flags)
