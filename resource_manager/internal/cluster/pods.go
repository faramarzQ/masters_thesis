package cluster

import (
	"context"
	"log"
	"resource_manager/config"

	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/metrics/pkg/apis/metrics/v1beta1"
)

type Pod struct {
	v1.Pod
	Name string
}

type PodList []Pod

func BindPod(pod v1.Pod) Pod {
	newPod := Pod{
		pod,
		pod.Name,
	}

	return newPod
}

// Returns all the metrics for the pod
func (p Pod) GetMetrics() *v1beta1.PodMetrics {
	podMetrics, err := MetricsClientset.MetricsV1beta1().PodMetricses(config.CLUSTER_NAMESPACE).Get(context.Background(), p.Name, metav1.GetOptions{})
	if err != nil {
		log.Fatal(err)
	}

	return podMetrics
}
