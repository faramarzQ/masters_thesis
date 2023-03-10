package cluster

import (
	"context"
	"fmt"
	"math"
	"scaler/config"
	"scaler/internal/consts"

	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/klog"
)

var (
	NodesClassCount = map[string]int{}
)

type Node struct {
	v1.Node
	Class                   string
	Hostname                string
	AllocatableCpu          string
	AllocatableMemory       string
	AllocatableStorage      string
	TotalCpu                int64
	TotalMemory             int64
	TotalStorage            string
	Architecture            string
	BootID                  string
	ContainerRuntimeVersion string
	MachineID               string
	SystemUUID              string
}

type NodeList []Node

func bindNode(node v1.Node) Node {
	totalCpu, _ := node.Status.Capacity.Cpu().AsInt64()
	totalMemory, _ := node.Status.Capacity.Memory().AsInt64()

	newNode := Node{
		node,
		node.Labels[consts.NODE_CLASS_LABEL_NAME],
		node.ObjectMeta.Name,
		node.Status.Allocatable.Cpu().String(),
		node.Status.Allocatable.Memory().String(),
		node.Status.Allocatable.StorageEphemeral().String(),
		totalCpu,
		totalMemory,
		node.Status.Capacity.StorageEphemeral().String(),
		node.Status.NodeInfo.Architecture,
		node.Status.NodeInfo.BootID,
		node.Status.NodeInfo.ContainerRuntimeVersion,
		node.Status.NodeInfo.MachineID,
		node.Status.NodeInfo.SystemUUID,
	}

	return newNode
}

func incrementNodeClassCount(class string) {
	NodesClassCount[class] += 1
}

func (n *Node) setClass(class string) {
	labelPatch := fmt.Sprintf(`[{"op":"add","path":"/metadata/labels/%s","value":"%s" }]`, "class", "off")
	_, err := Clientset.CoreV1().Nodes().Patch(context.Background(), n.Name, types.JSONPatchType, []byte(labelPatch), metav1.PatchOptions{})
	if err != nil {
		panic(err)
	}
}

func resetNodesClassCountToZero() {
	NodesClassCount = map[string]int{
		consts.ACTIVE_CLASS: 0,
		consts.IDLE_CLASS:   0,
		consts.SLEEP_CLASS:  0,
		consts.OFF_CLASS:    0,
	}
}

func (n Node) ListPods() PodList {
	pods, err := Clientset.CoreV1().Pods(config.CLUSTER_NAMESPACE).List(context.Background(), metav1.ListOptions{
		FieldSelector: "spec.nodeName=" + n.Name,
	})
	if err != nil {
		panic(err)
	}

	podList := PodList{}
	for _, pod := range pods.Items {
		podList = append(podList, bindPod(pod))
	}

	return podList
}

// Calculate node's memory utilization
func (n Node) GetMemoryUtilization() float64 {
	totalMemory := n.TotalMemory
	var usedMemory int64

	for _, pod := range n.ListPods() {
		podMetrics := pod.GetMetrics()
		for _, container := range podMetrics.Containers {
			containerMemoryUsage, _ := container.Usage.Memory().AsInt64()
			usedMemory += containerMemoryUsage
		}
	}
	memoryUtilization := (float64(usedMemory) / float64(totalMemory)) * 100
	memoryUtilizationRounded := math.Floor(memoryUtilization*10000) / 10000

	return memoryUtilizationRounded
}

// Calculate node's cpu utilization
func (n Node) GetCpuUtilization() float64 {
	totalCpu := n.TotalCpu
	var usedCpu int64

	for _, pod := range n.ListPods() {
		podMetrics := pod.GetMetrics()
		for _, container := range podMetrics.Containers {
			containerCpuUsage, _ := container.Usage.Cpu().AsInt64()
			usedCpu += containerCpuUsage
		}
	}
	cpuUtilization := (float64(usedCpu) / float64(totalCpu)) * 100
	cpuUtilizationRounded := math.Floor(cpuUtilization*10000) / 10000

	return cpuUtilizationRounded
}

// List all nodes in the cluster
func ListNodes() NodeList {
	nodes, err := Clientset.CoreV1().Nodes().List(context.Background(), metav1.ListOptions{})
	if err != nil {
		klog.Fatal(err)
	}

	nodeList := NodeList{}
	// resetNodesClassCountToZero()
	for _, node := range nodes.Items {
		wrappedNode := bindNode(node)
		nodeList = append(nodeList, wrappedNode)
		// incrementNodeClassCount(wrappedNode.Class)
	}

	return nodeList
}

// List all the active nodes in the cluster
func ListActiveNodes() NodeList {
	nodes := ListNodes()
	nodeList := NodeList{}
	for _, node := range nodes {
		if node.Class == consts.ACTIVE_CLASS {
			nodeList = append(nodeList, node)
		}
	}

	return nodeList
}

// add off label to nodes without class label
func RegisterNodes() {
	nodeList := ListNodes()
	for _, node := range nodeList {
		if node.Class == "" {
			node.setClass(consts.OFF_CLASS)
		}
	}
}
