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

RANDOM_SCHEDULER_DEFAULT_FLAGS=--authentication-kubeconfig=/etc/kubernetes/scheduler.conf --authorization-kubeconfig=/etc/kubernetes/scheduler.conf --config=../config/scheduler/random-scheduler-config.yml --secure-port=10260
HEURISTIC_SCHEDULER_DEFAULT_FLAGS=--authentication-kubeconfig=/etc/kubernetes/scheduler.conf --authorization-kubeconfig=/etc/kubernetes/scheduler.conf --config=../config/scheduler/heuristic-scheduler-config.yml --secure-port=10260
DEFAULT_SCHEDULER_DEFAULT_FLAGS=--authentication-kubeconfig=/etc/kubernetes/scheduler.conf --authorization-kubeconfig=/etc/kubernetes/scheduler.conf --config=../config/scheduler/default-scheduler-config.yml --secure-port=10260

SCALER_IMAGE_NAME=scaler
SCHEDULER_IMAGE_NAME=scheduler

SCALER_DOCKERFILE=scaler.Dockerfile
SCHEDULER_DOCKERFILE=scheduler.Dockerfile

#
# Scaler
#

# On local

run-random-scaler:
	cd $(SCRIPTS_DIR) && . ./$(ENV_LOCAL) \
	&& cd $(PROJECT_ROOT) && APP=scaler APP_ID=1 go run $(SCALER_ENTRY) $(flags)

run-fixed-scaler:
	cd $(SCRIPTS_DIR) && . ./$(ENV_LOCAL) \
	&& cd $(PROJECT_ROOT) && APP=scaler APP_ID=2 go run $(SCALER_ENTRY) $(flags)

run-heuristic-scaler:
	cd $(SCRIPTS_DIR) && . ./$(ENV_LOCAL) \
	&& cd $(PROJECT_ROOT) && APP=scaler APP_ID=3 go run $(SCALER_ENTRY) $(flags)

run-proposed-scaler:
	cd $(SCRIPTS_DIR) && . ./$(ENV_LOCAL) \
	&& cd $(PROJECT_ROOT) && APP=scaler APP_ID=4 go run $(SCALER_ENTRY) $(flags)

run-silencer-scaler:
	cd $(SCRIPTS_DIR) && . ./$(ENV_LOCAL) \
	&& cd $(PROJECT_ROOT) && APP=scaler APP_ID=5 go run $(SCALER_ENTRY) $(flags)

build-scaler-image:
	docker build --tag $(SCALER_IMAGE_NAME):1.0 . -f $(SCALER_DOCKERFILE)

# On Minikube

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

run-proposed-scaler-bin:
	cd $(SCRIPTS_DIR) && . ./$(ENV_MINIKUBE_BIN) \
	&& cd $(BIN_DIR) && APP=scaler APP_ID=4 ./scaler $(flags)

run-silencer-scaler-bin:
	cd $(SCRIPTS_DIR) && . ./$(ENV_MINIKUBE_BIN) \
	&& cd $(BIN_DIR) && APP=scaler APP_ID=5 ./scaler $(flags)

#
# Scheduler
#

run-random-scheduler:
	cd $(SCRIPTS_DIR) && . ./$(ENV_LOCAL) \
	&& cd $(PROJECT_ROOT) && APP=scheduler APP_ID=1 go run $(SCHEDULER_ENTRY) $(flags)

build-scheduler:
	go build -o $(SCHEDULER_OUTPUT) $(SCHEDULER_ENTRY) $(flags)

run-random-scheduler-bin:
	cd $(SCRIPTS_DIR) && . ./$(ENV_MINIKUBE_BIN) \
	&& cd $(BIN_DIR) && APP=scheduler APP_ID=1 ./scheduler $(RANDOM_SCHEDULER_DEFAULT_FLAGS) $(flags)

run-heuristic-scheduler-bin:
	cd $(SCRIPTS_DIR) && . ./$(ENV_MINIKUBE_BIN) \
	&& cd $(BIN_DIR) && APP=scheduler APP_ID=2 ./scheduler $(HEURISTIC_SCHEDULER_DEFAULT_FLAGS) $(flags)

run-default-scheduler-bin:
	cd $(SCRIPTS_DIR) && . ./$(ENV_MINIKUBE_BIN) \
	&& cd $(BIN_DIR) && APP=scheduler APP_ID=3 ./scheduler $(DEFAULT_SCHEDULER_DEFAULT_FLAGS) $(flags)

build-scheduler-image:
	docker build --tag $(SCHEDULER_IMAGE_NAME):1.0 . -f $(SCHEDULER_DOCKERFILE)
