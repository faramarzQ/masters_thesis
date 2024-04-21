# masters_thesis
My master's thesis project implementing a resource management service in serverless-edge architecture

# Setup
    - install minikube
        https://minikube.sigs.k8s.io/docs/start/

    - install kubectl
        https://kubernetes.io/docs/tasks/tools/install-kubectl-linux/\

    - install open-faas
        https://hackernoon.com/how-to-create-serverless-functions-with-openfaas-in-17-steps-u21l3y7m

    - make access to openfaas panel
        `minikube -p thesis-cluster -n openfaas service --url gateway-external`

    - create nodeport service for prometheus and give access using minikube:
        `kubectl expose deployment prometheus -n openfaas --type=NodePort --name=prometheus-external`
        `minikube -p thesis-cluster -n openfaas service --url prometheus-external`

    - download Grafana image:
        `kubectl run grafana -n openfaas --image=stefanprodan/faas-grafana:4.6.3 --port=3000`
        `kubectl expose pod grafana -n openfaas --type=NodePort --name=grafana`
        `minikube -p thesis-cluster -n openfaas service --url grafana`

    - expose prometheus' dashboard using: 
        `minikube -p thesis-cluster -n openfaas service wprometheus-external`



-----------------

expose a specific pod:

kubectl expose pod fibonacci-6cc54fb8fb-559tw  -n openfaas-fn --type=NodePort --name=grafana

minikube -p thesis-cluster -n openfaas-fn service grafana

-----------------

setup prometheus:

https://devopscube.com/setup-prometheus-monitoring-on-kubernetes/


-----------------

# Projects explanation:

## Monitoring server

This service is responsible for sending the application's metrics to Prometheus by Exposing HTTP APIs.

## Gateway

Gateway service is a middleware layer between the client and the handler function, responsible for:
- Error handling including timeout errors
- Sending request metrics to Prometheus including:
    - response time
    - incrementing total number of requests
    - incrementing total number of failed requests

## Fibonacci Function

Fibonacci function is the request handler of the application, processing the fibonacci number of a given number received from exposed HTTP API. This function applies timeout on the calculation of the numbers and responds with error if timeout passed. in case of success calculation of the number, it also increments the success requests of Prometheus

## Resource Manager

Resource manager contains is responsible for managing and controlling the usage of the cluster's resource, including the scaling of the nodes, scheduling of the functions, etc.

- ### Scaler Manager
The scaler manager application registers and manages the execution of different scaler applications including:     
    - **Random Scaler**: It scales a specific number of nodes between off and idle groups by based on a probability number.    
    - **Fixed Scaler**: It keeps the number of Idle nodes at a specific number, by turning on Off nodes.    
    - **Heuristic Scaler**: Scales the nodes based on the active nodes' CPU utilization, If passed a specific threshold, adds enough nodes from off to idle for it to reach a desired threshold. Most efficient nodes are being selected to be scaled.    
    - **Proposed Scaler**: The proposed scaler communicates with the AI agent to receive the number of nodes to be scaled between idle and off classes, By sending the current status of the cluster including the request rate, energy consumption, etc. the AI agent responds with the number of nodes to be transitioned between off and idle classes. Most efficient nodes are being selected to be transitioned.    
    - **Silencer scaler**: Silencer scaler is rather a different kind of scaler with should be executed along other scalers. It scales down nodes to lower classes and is being executed in smaller period times.    

- ### Scheduler Manager
The scheduler manager application registers and manages the execution of different scheduler application including:     
    - **Default Scheduler**: Implements the default scheduler of Kubernetes which schedules the pods on nodes based on sparse algorithm.    
    - **Random Scheduler**: Schedules pods on random nodes.    
    - **Heuristic Scheduler**: Schedules pods on nodes with most cpu utilization.    

## AI Agent
The AI agent is an AI application is being used by the proposed scaler to predict the best number of nodes to be scaled between different classes. It exposes and API to receive the cluster's status and responds the proper action. It uses a reinforcement learning model to predict the actions.



-------------------

# Running the project

- ## Scaler
    - Apply gateway and fibonacci deployments
    - Apply HPA for fibonacci
    - Add these anottation to the nodes:
    >    max_power_consumption: '100'
    >    min_power_consumption: '10'
    - mount the resource manager's directory to minikube using this command:    
    > minikube -p thesis-cluster mount ./:/home/telepathy
    - build the resource manager's binary on the local using:
    > make build-scaler
    - expose cluster's database using:
    > m service database-nodeport
    - login to minikube's ssh using:
    - minikube -p thesis-cluster ssh
    - run scaler using:
    > DATABASE_NODEPORT_SERVICE_HOST=192.168.49.2 DATABASE_NODEPORT_SERVICE_PORT=30306 make run-random-scaler-bin