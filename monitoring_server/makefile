PROJECT_ROOT := $(realpath $(dir $(firstword $(MAKEFILE_LIST))))

IMAGE_NAME=monitoring_server

run-monitoring-server:
	python3 main.py

build-monitoring-server:
	docker build --tag $(IMAGE_NAME):1.0 .