#!/bin/bash

export ENV=DEV_LOCAL
export CONFIG_DIR_DEV_LOCAL=../config/outcluster_auth/config.yml
export CLUSTER_NAMESPACE=default

# Database
export DB_HOST=0.0.0.0
export DB_PORT=3306
export DB_DATABASE=cluster
export DB_USERNAME=sample_user
export DB_PASSWORD=9xz3jrd8wf

export AI_SERVER_URL=http://localhost:8080

# Prometheus
export PROMETHEUS_URL=http://192.168.49.2:31090
export PROMETHEUS_SUCCESS_REQUESTS_PERIOD_MINUTE=0
export PROMETHEUS_METRIC_NAME_SUCCESS_REQUESTS=success_requests_total
export PROMETHEUS_TOTAL_REQUESTS_PERIOD_MINUTE=0
export PROMETHEUS_METRIC_NAME_TOTAL_REQUESTS=requests_total

# RL Agent
export RL_SUCCESS_RATE_WEIGHT=2
export RL_ENERGY_CONSUMPTION_WEIGHT=3
export RL_ALFA_VALUE="0.5"
export RL_GAMMA_VALUE="0.6"
export RL_MAXIMUM_EPSILON_VALUE=0.8
export RL_MINIMUM_EPSILON_VALUE=0.05
export RL_EDR=0.02

#  Fixed scaler
export FIXED_IDLE_NODES_COUNT=2
