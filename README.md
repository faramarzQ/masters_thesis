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