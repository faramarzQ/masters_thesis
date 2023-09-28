# App env vars
FIBONACCI_TIMEOUT_SECONDS=10

IMAGE_NAME=fibonacci

run-fibonacci:
	FIBONACCI_TIMEOUT_SECONDS=$(FIBONACCI_TIMEOUT_SECONDS) \
	go run main.go

build-fibonacci-image:
	docker build --tag $(IMAGE_NAME):1.0 .