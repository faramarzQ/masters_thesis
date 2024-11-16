package cluster

import (
	"context"
	"fmt"
	"log"
	"resource_manager/internal/config"
	"resource_manager/internal/consts"
	"time"

	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/klog"
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

// Syncs pod with the cluster
func (p *Pod) Update() {
	newPod, err := Clientset.CoreV1().Pods(config.CLUSTER_NAMESPACE).Get(context.Background(), p.Name, metav1.GetOptions{})
	if err != nil {
		klog.Fatal(err)
	}

	wrappedPod := BindPod(*newPod)
	*p = wrappedPod
}

// Sets an annotation on the pod
func (p *Pod) SetAnnotation(key, value string) {
	p.SetAnnotations(map[string]string{key: value})
	_, err := Clientset.CoreV1().Pods(config.CLUSTER_NAMESPACE).Update(context.TODO(), &p.Pod, metav1.UpdateOptions{})
	if err != nil {
		panic(err)
	}
}

// Gets annotation value for a given key on the pod
func (p *Pod) GetAnnotation(key string) string {
	pod, _ := Clientset.CoreV1().Pods(config.CLUSTER_NAMESPACE).Get(context.TODO(), p.Name, metav1.GetOptions{})
	for annotation_name, annotation_value := range pod.GetAnnotations() {
		if annotation_name == key {
			return annotation_value
		}
	}

	return ""
}

// Updates the node's class label
func (p *Pod) WarmUp() {
	labelPatch := fmt.Sprintf(`[{"op":"add","path":"/metadata/labels/%s","value":"%s" }]`, consts.POD_WARM_LABEL_NAME, "true")
	_, err := Clientset.CoreV1().Pods(config.CLUSTER_NAMESPACE).Patch(context.Background(), p.Name, types.JSONPatchType, []byte(labelPatch), metav1.PatchOptions{})
	if err != nil {
		panic(err)
	}

	p.Update()
	p.SetWarmedAt(time.Now())
}

// Updates the pod's warmed_at timestamp
func (p *Pod) SetWarmedAt(timestamp time.Time) {
	formatted := timestamp.Format(time.RFC3339)
	p.SetAnnotation(consts.POD_WARMED_AT_ANNOTATION_NAME, formatted)
}

// Gets Warmed_at label on pod
func (p *Pod) GetWarmedAt() time.Time {
	warmedAt := p.GetAnnotation(consts.POD_WARMED_AT_ANNOTATION_NAME)
	date, error := time.Parse(time.RFC3339, warmedAt)
	if error != nil {
		panic(error)
	}

	return date
}

// Checks if pod is labeled as warm
func (p *Pod) IsAlreadyWarm() bool {
	isWarm := p.Labels[consts.POD_WARM_LABEL_NAME]
	return isWarm == "true"
}

// Set warm label as false
func (p *Pod) UnsetWarm() {
	labelPatch := fmt.Sprintf(`[{"op":"add","path":"/metadata/labels/%s","value":"%s" }]`, consts.POD_WARM_LABEL_NAME, "false")
	_, err := Clientset.CoreV1().Pods(config.CLUSTER_NAMESPACE).Patch(context.Background(), p.Name, types.JSONPatchType, []byte(labelPatch), metav1.PatchOptions{})
	if err != nil {
		panic(err)
	}
}

// Returns all the metrics for the pod
func (p Pod) GetMetrics() *v1beta1.PodMetrics {
	podMetrics, err := MetricsClientset.MetricsV1beta1().PodMetricses(config.CLUSTER_NAMESPACE).Get(context.Background(), p.Name, metav1.GetOptions{})
	if err != nil {
		log.Fatal(err)
	}

	return podMetrics
}

func ListAllPods() PodList {
	pods, err := Clientset.CoreV1().Pods(config.CLUSTER_NAMESPACE).List(context.Background(), metav1.ListOptions{
		LabelSelector: "app=fibonacci",
	})
	if err != nil {
		panic(err)
	}

	podList := PodList{}
	for _, pod := range pods.Items {
		podList = append(podList, BindPod(pod))
	}

	return podList
}
