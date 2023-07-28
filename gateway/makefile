PROJECT_ROOT := $(realpath $(dir $(firstword $(MAKEFILE_LIST))))
ENV_LOCAL := env_local.sh

run-gateway:
	. ./$(ENV_LOCAL) \
	&& go run main.go
