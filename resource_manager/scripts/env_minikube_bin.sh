#!/bin/bash

# This config is for debuggin from minikube sh by executing binary files

export ENV=DEV_LOCAL
export CONFIG_DIR_DEV_LOCAL=../config/outcluster_auth/config.yml
export CLUSTER_NAMESPACE=default
export LOG_FILE_DIR=storage/logs.log
export LOG_TO_CONSOLE="true"

# Database
export DB_HOST=192.168.49.2 
export DB_PORT=30306
export DB_DATABASE=cluster
export DB_USERNAME=sample_user
export DB_PASSWORD=9xz3jrd8wf

# Prometheus
export PROMETHEUS_URL=http://192.168.49.2:31090
export PROMETHEUS_SUCCESS_REQUESTS_PERIOD_MINUTE=0
export PROMETHEUS_METRIC_NAME_SUCCESS_REQUESTS=success_requests_total
export PROMETHEUS_TOTAL_REQUESTS_PERIOD_MINUTE=0
export PROMETHEUS_METRIC_NAME_TOTAL_REQUESTS=requests_total

# RL Agent
export AI_AGENT_URL=http://192.168.49.2:31000/
export RL_SUCCESS_RATE_WEIGHT=2
export RL_ENERGY_CONSUMPTION_WEIGHT=3
export RL_ALFA_VALUE="0.5"
export RL_GAMMA_VALUE="0.6"
export RL_MAXIMUM_EPSILON_VALUE=0.99
export RL_MINIMUM_EPSILON_VALUE=0.05
export RL_EDR=0.02

# Fixed scaler
export FIXED_IDLE_NODES_COUNT=2

# Monitoring server
export MONITORING_SERVER_URL=http://192.168.49.2:32000/

# Heuristic scaler
export HEURISTIC_SCALER_UPPER_CPU_THRESHOLD=10
export HEURISTIC_SCALER_DESIRED_CPU_UTIL=5

# Silencer scaler
export SILENCER_SCALER_WARM_POD_DURATION_MINUTES=1
export SILENCER_SCALER_IDLE_NODE_DURATION_MINUTES=10

# Random Scaler
export RANDOM_SCALER_SCALING_PROBABILITY=60
export RANDOM_SCALER_PERCENT_OF_NODES_TO_TRANSIT=30
