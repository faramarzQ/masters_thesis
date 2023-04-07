PROJECT_ROOT := $(realpath $(dir $(firstword $(MAKEFILE_LIST))))
SCALER_ENTRY := $(PROJECT_ROOT)/cmd/scaler/main.go

run-scaler:
	. ./env.sh && ENV=DEV_LOCAL go run $(SCALER_ENTRY) $(flags)