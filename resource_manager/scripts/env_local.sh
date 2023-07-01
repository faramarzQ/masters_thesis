#!/bin/bash

export ENV=DEV_LOCAL
export CONFIG_DIR_DEV_LOCAL=../config/outcluster_auth/config.yml
export CLUSTER_NAMESPACE=default

export DB_HOST=0.0.0.0
export DB_PORT=3306
export DB_DATABASE=cluster
export DB_USERNAME=sample_user
export DB_PASSWORD=9xz3jrd8wf

export AI_SERVER_URL=http://localhost:8080

export PROMETHEUS_URL=http://192.168.49.2:31090
export PROMETHEUS_SUCCESS_REQUESTS_PERIOD_MINUTE=0
export PROMETHEUS_METRIC_NAME_SUCCESS_REQUESTS=success_requests_total