PROJECT_ROOT := $(realpath $(dir $(firstword $(MAKEFILE_LIST))))
SILENCER_ENTRY := $(PROJECT_ROOT)/cmd/silencer/main.go
AI_AGENT := $(PROJECT_ROOT)/ai_agent
DATABASE := $(PROJECT_ROOT)/database
FIBONACCI := $(PROJECT_ROOT)/fibonacci
GATEWAY := $(PROJECT_ROOT)/gateway
GRAFANA := $(PROJECT_ROOT)/grafana
MONITORING_SERVER := $(PROJECT_ROOT)/monitoring_server
NODE_EXPORTER := $(PROJECT_ROOT)/node_exporter
PROMETHEUS := $(PROJECT_ROOT)/prometheus
FIXED_SCALER := $(PROJECT_ROOT)/resource_manager/fixed_scaler
HEURISTIC_SCALER := $(PROJECT_ROOT)/resource_manager/heuristic_scaler
PROPOSED_SCALER := $(PROJECT_ROOT)/resource_manager/proposed_scaler
RANDOM_SCALER := $(PROJECT_ROOT)/resource_manager/random_scaler
SILENCER_SCALER := $(PROJECT_ROOT)/resource_manager/silencer_scaler

KAF := kubectl apply -f

deploy-ai-agent:
	cd $(AI_AGENT) && \
	$(KAF) deployment.yml && \
	$(KAF) nodePort.yml

deploy-database:
	cd $(DATABASE) && \
	$(KAF) nodePort.yml && \
	$(KAF) configMap.yml && \
	$(KAF) replicaSet.yml

deploy-fibonacci:
	cd $(FIBONACCI) && \
	$(KAF) nodePort.yml && \
	$(KAF) deployment.yml && \
	$(KAF) horizontalPodAutoscaler.yml

deploy-gateway:
	cd $(GATEWAY) && \
	$(KAF) nodePort.yml && \
	$(KAF) deployment.yml && \
	$(KAF) horizontalPodAutoscaler.yml

deploy-grafana:
	cd $(GRAFANA) && \
	$(KAF) nodePort.yml && \
	$(KAF) replicaSet.yml

deploy-monitoring-server:
	cd $(GATEWAY) && \
	$(KAF) nodePort.yml && \
	$(KAF) deployment.yml

deploy-node-exporter:
	cd $(NODE_EXPORTER) && \
	$(KAF) daemonSet.yml

deploy-prometheus:
	cd $(PROMETHEUS) && \
	$(KAF) clusterRole.yml && \
	$(KAF) serviceAccount.yml && \
	$(KAF) nodePort.yml && \
	$(KAF) clusterIp.yml && \
	$(KAF) configMap.yml && \
	$(KAF) nodePort.yml

deploy-fixed-scaler:
	cd $(FIXED_SCALER) && \
	$(KAF) cronJob.yml

deploy-heuristic-scaler:
	cd $(HEURISTIC_SCALER) && \
	$(KAF) cronJob.yml

deploy-proposed-scaler:
	cd $(PROPOSED_SCALER) && \
	$(KAF) cronJob.yml

deploy-random-scaler:
	cd $(RANDOM_SCALER) && \
	$(KAF) cronJob.yml

deploy-silencer-scaler:
	cd $(SILENCER_SCALER) && \
	$(KAF) cronJob.yml
