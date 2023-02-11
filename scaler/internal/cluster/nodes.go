package cluster

import (
	"context"
	"fmt"
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
	*v1.Node
	Class                   string
	Hostname                string
	AllocatableCpu          string
	AllocatableMemory       string
	AllocatableStorage      string
	TotalCpu                string
	TotalMemory             string
	TotalStorage            string
	Architecture            string
	BootID                  string
	ContainerRuntimeVersion string
	MachineID               string
	SystemUUID              string
}

type NodeList []*Node

func bindNode(node *v1.Node) *Node {
	newNode := &Node{
		node,
		node.Labels[consts.NODE_CLASS_LABEL_NAME],
		node.ObjectMeta.Name,
		node.Status.Allocatable.Cpu().String(),
		node.Status.Allocatable.Memory().String(),
		node.Status.Allocatable.StorageEphemeral().String(),
		node.Status.Capacity.Memory().String(),
		node.Status.Capacity.StorageEphemeral().String(),
		node.Status.Capacity.Cpu().String(),
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
	labelPatch := fmt.Sprintf(`[{"op":"add","path":"/metadata/labels/%s","value":"%s" }]`, consts.NODE_CLASS_LABEL_NAME, class)
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

func ListNodes() NodeList {
	nodes, err := Clientset.CoreV1().Nodes().List(context.Background(), metav1.ListOptions{})
	if err != nil {
		klog.Fatal(err)
	}

	nodeList := NodeList{}
	resetNodesClassCountToZero()
	for _, node := range nodes.Items {
		wrappedNode := bindNode(&node)
		nodeList = append(nodeList, wrappedNode)
		incrementNodeClassCount(wrappedNode.Class)
	}

	fmt.Println(NodesClassCount)

	return nodeList
}

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

func RegisterNodes() {
	nodeList := ListNodes()
	for _, node := range nodeList {
		if node.Class == "" {
			node.setClass(consts.OFF_CLASS)
		}
	}
}
