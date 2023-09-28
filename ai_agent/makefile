PROJECT_ROOT := $(realpath $(dir $(firstword $(MAKEFILE_LIST))))

IMAGE_NAME=ai_agent

run-agent:
	python3 main.py

build-ai-agent-image:
	docker build --tag $(IMAGE_NAME):1.0 .