#!/bin/bash

export ENV=DEV_LOCAL
export CONFIG_DIR_DEV_LOCAL=../config/outcluster_auth/config.yml
export CLUSTER_NAMESPACE=openfaas-fn

export DB_HOST=0.0.0.0
export DB_PORT=3306
export DB_DATABASE=cluster
export DB_USERNAME=sample_user
export DB_PASSWORD=9xz3jrd8wf