package cluster

import (
	"context"
	"errors"
	"fmt"
	"math"
	"math/rand"
	"resource_manager/internal/config"
	"resource_manager/internal/consts"
	"resource_manager/internal/helpers"
	"resource_manager/internal/prometheus"
	"strconv"
	"time"

	"github.com/prometheus/common/model"
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
	Class                   consts.NODE_CLASS
	IsMaster                bool
	IsWorker                bool
	Hostname                string
	AllocatableCpu          string
	AllocatableMemory       string
	AllocatableStorage      string
	TotalCpu                int64
	TotalMemory             int64
	TotalStorage            string
	MaxPowerConsumption     int64
	MinPowerConsumption     int64
	Architecture            string
	BootID                  string
	ContainerRuntimeVersion string
	MachineID               string
	SystemUUID              string
}

type NodeList []Node

// Returns string names of the node list
func (nl NodeList) Names() []string {
	var names []string
	for _, node := range nl {
		names = append(names, node.Name)
	}
	return names
}

// Returns nodes from the node list which are in the given classes
func (nl NodeList) InClass(classes ...consts.NODE_CLASS) NodeList {
	var newNodeList NodeList
	for _, node := range nl {
		for _, class := range classes {
			if node.Class == class {
				newNodeList = append(newNodeList, node)
			}
		}
	}
	return newNodeList
}

// Return total memory amount on the given node list
func (nl NodeList) TotalMemory() uint {
	var totalMemory uint
	for _, node := range nl {
		totalMemory += uint(node.TotalMemory)
	}
	return totalMemory
}

// Return total cpu amount on the given node list
func (nl NodeList) TotalCpu() uint {
	var totalCpu uint
	for _, node := range nl {
		totalCpu += uint(node.TotalCpu)
	}
	return totalCpu
}

func BindNode(node v1.Node) Node {
	totalCpu, _ := node.Status.Capacity.Cpu().AsInt64() // ALERT: don't use this! use n.GetTotalCpuCores() instead.
	totalMemory, _ := node.Status.Capacity.Memory().AsInt64()
	var class consts.NODE_CLASS = consts.NODE_CLASS(node.Labels[consts.NODE_CLASS_LABEL_NAME])
	isMaster, _ := strconv.ParseBool(node.Labels[consts.NODE_IS_PRIMARY_LABEL_NAME])
	isWorker, _ := strconv.ParseBool(node.Labels[consts.WORKER_NODE_LABEL_NAME])

	var maxPowerConsumption, minPowerConsumption int
	var err error
	if isWorker {
		maxPowerConsumption, err = strconv.Atoi(node.Annotations[consts.MAX_POWER_CONSUMPTION_LABEL_NAME])
		if err != nil {
			klog.Fatal(err)
		}

		minPowerConsumption, err = strconv.Atoi(node.Annotations[consts.MIN_POWER_CONSUMPTION_LABEL_NAME])
		if err != nil {
			klog.Fatal(err)
		}
	}

	newNode := Node{
		node,
		class,
		isMaster,
		isWorker,
		node.ObjectMeta.Name,
		node.Status.Allocatable.Cpu().String(),
		node.Status.Allocatable.Memory().String(),
		node.Status.Allocatable.StorageEphemeral().String(),
		totalCpu,
		totalMemory,
		node.Status.Capacity.StorageEphemeral().String(),
		int64(maxPowerConsumption),
		int64(minPowerConsumption),
		node.Status.NodeInfo.Architecture,
		node.Status.NodeInfo.BootID,
		node.Status.NodeInfo.ContainerRuntimeVersion,
		node.Status.NodeInfo.MachineID,
		node.Status.NodeInfo.SystemUUID,
	}

	return newNode
}

// Syncs node with the cluster
func (n *Node) Update() {
	newNode, err := Clientset.CoreV1().Nodes().Get(context.Background(), n.Name, metav1.GetOptions{})
	if err != nil {
		klog.Fatal(err)
	}

	wrappedNode := BindNode(*newNode)
	*n = wrappedNode
}

// Updates the node's class label
func (n *Node) SetClass(class consts.NODE_CLASS) {
	labelPatch := fmt.Sprintf(`[{"op":"add","path":"/metadata/labels/%s","value":"%s" }]`, consts.NODE_CLASS_LABEL_NAME, class)
	newNode, err := Clientset.CoreV1().Nodes().Patch(context.Background(), n.Name, types.JSONPatchType, []byte(labelPatch), metav1.PatchOptions{})
	if err != nil {
		panic(err)
	}

	wrappedNode := BindNode(*newNode)
	*n = wrappedNode
	n.SetScaledAt(time.Now())
}

// Sets an annotation on the node
func (n *Node) SetAnnotation(key, value string) {
	annotations := n.Annotations
	annotations[key] = value
	n.SetAnnotations(annotations)
	newNode, err := Clientset.CoreV1().Nodes().Update(context.TODO(), &n.Node, metav1.UpdateOptions{})
	if err != nil {
		panic(err)
	}

	wrappedNode := BindNode(*newNode)
	*n = wrappedNode
}

// Gets annotation value for a given key on the node
func (n *Node) GetAnnotation(key string) string {
	node, _ := Clientset.CoreV1().Nodes().Get(context.TODO(), n.Name, metav1.GetOptions{})
	for annotation_name, annotation_value := range node.GetAnnotations() {
		if annotation_name == key {
			return annotation_value
		}
	}

	return ""
}

// Updates the node's scaled_at timestamp
func (n *Node) SetScaledAt(timestamp time.Time) {
	formatted := timestamp.Format(time.RFC3339)
	n.SetAnnotation(consts.NODE_SCALED_AT_LABEL_NAME, formatted)
}

// Gets scaled_at label on pod
func (n *Node) GetScaledAt() time.Time {
	scaledAt := n.GetAnnotation(consts.NODE_SCALED_AT_LABEL_NAME)
	date, error := time.Parse(time.RFC3339, scaledAt)
	if error != nil {
		panic(error)
	}

	return date
}

func (n Node) ListPods() PodList {
	pods, err := Clientset.CoreV1().Pods(config.CLUSTER_NAMESPACE).List(context.Background(), metav1.ListOptions{
		FieldSelector: "spec.nodeName=" + n.Name,
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
	totalCpu := n.GetTotalCpuCores()
	usedCpu, err := n.GetUsedCpuCoresAtGiveTime(time.Now())
	if err != nil {
		klog.Fatal(usedCpu)
	}

	cpuUtilization := (float64(usedCpu) / float64(totalCpu)) * 100
	cpuUtilizationRounded := math.Floor(cpuUtilization*1) / 1

	return cpuUtilizationRounded
}

// Returns the efficiency of the node based on it's total memory and power usage,
// the higher the better
func (n Node) GetMemoryEfficiency() int64 {
	efficiency := n.TotalMemory / n.MaxPowerConsumption // TODO: should be max?
	return efficiency
}

// Returns the efficiency of the node based on it's total cpu and power usage,
// the higher the better
func (n Node) GetCpuEfficiency() float64 {
	efficiency := float64(n.TotalCpu) / float64(n.MaxPowerConsumption) // TODO: should be max?
	return efficiency
}

// Return actual number of cores of the node
func (n Node) GetTotalCpuCores() float64 {
	result := prometheus.Query("machine_cpu_physical_cores{instance='"+n.Name+"'}", time.Now())
	return float64(result.(model.Vector)[0].Value)
}

func (n Node) GetUsedCpuCoresAtGiveTime(time time.Time) (float64, error) {
	result := prometheus.Query("sum(rate(container_cpu_usage_seconds_total{instance=~'"+n.Name+"', pod=~'fibonacci.*'}[1m]))", time)
	if len(result.(model.Vector)) != 0 {
		return float64(result.(model.Vector)[0].Value), nil
	}
	return 0, errors.New("Failed finding any data at the given time")
}

// List all nodes in the cluster
func ListNodes() NodeList {
	listOptions := metav1.ListOptions{
		LabelSelector: consts.WORKER_NODE_LABEL_NAME + "=true",
	}
	nodes, err := Clientset.CoreV1().Nodes().List(context.Background(), listOptions)
	if err != nil {
		klog.Fatal(err)
	}

	nodeList := NodeList{}
	for _, node := range nodes.Items {
		wrappedNode := BindNode(node)
		nodeList = append(nodeList, wrappedNode)
	}

	return nodeList
}

// List all nodes in the cluster
func ListAllNodes() NodeList {
	listOptions := metav1.ListOptions{}
	nodes, err := Clientset.CoreV1().Nodes().List(context.Background(), listOptions)
	if err != nil {
		klog.Fatal(err)
	}

	nodeList := NodeList{}
	for _, node := range nodes.Items {
		wrappedNode := BindNode(node)
		nodeList = append(nodeList, wrappedNode)
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
func LabelNewNodes() {
	nodeList := ListNodes()
	for _, node := range nodeList {
		if node.Class == "" {
			node.SetClass(consts.OFF_CLASS)
		}
	}
}

// Returns the master node
func MasterNode() *Node {
	for _, node := range ListAllNodes() {
		if node.IsMaster == true {
			return &node
		}
	}

	return nil
}

// Returns dispersion of nodes across all functioning classes
func getsNodesDispersion() map[consts.NODE_CLASS]int {
	nodes := ListNodes()

	// init state
	state := make(map[consts.NODE_CLASS]int)
	for _, class := range consts.FUNCTIONING_CLASSES {
		state[class] = 0
	}

	for _, node := range nodes {
		state[node.Class] += 1
	}

	return state
}

// Returns a randomly selected node from a list of nodes
func GetRandomNodesFromNodeList(nodeList NodeList, numberOfNodesToSelect int64) NodeList {
	nodesIndexesToSelect := []int{}

	for len(nodesIndexesToSelect) < int(numberOfNodesToSelect) {
		rand.Seed(time.Now().UnixNano())
		randomNum := rand.Intn(int(len(nodeList)))

		var nodeAlreadySelected bool
		for i := 0; i < len(nodesIndexesToSelect); i++ {
			if randomNum == nodesIndexesToSelect[i] {
				nodeAlreadySelected = true
				break
			}
		}

		if !nodeAlreadySelected {
			nodesIndexesToSelect = append(nodesIndexesToSelect, randomNum)
		}
	}

	var nodesToSelect NodeList
	for i := 0; i < len(nodesIndexesToSelect); i++ {
		nodesToSelect = append(nodesToSelect, nodeList[nodesIndexesToSelect[i]])
	}

	return nodesToSelect
}

// Finds the most memory efficient node in off class
// Returns a node and a flag indicating if any node found
func GetMostMemoryEfficientNode(exceptionNodesName []string, className consts.NODE_CLASS) (*Node, bool) {
	nodes := ListNodes().InClass(className)

	var maxEfficiency int64
	var mostEfficientNode Node
	for _, node := range nodes {
		temp := node.GetMemoryEfficiency()
		if temp > maxEfficiency && !helpers.StringInSlice(node.Name, exceptionNodesName) {
			maxEfficiency = temp
			mostEfficientNode = node
		}
	}

	if maxEfficiency == 0 {
		return &mostEfficientNode, false
	}

	return &mostEfficientNode, true
}

// Finds the most CPU efficient node in off class
// Returns a node and a flag indicating if any node found
func GetMostCpuEfficientNode(exceptionNodesName []string, className consts.NODE_CLASS) (*Node, bool) {
	nodes := ListNodes().InClass(className)

	var maxEfficiency float64
	var mostEfficientNode Node
	for _, node := range nodes {
		temp := node.GetCpuEfficiency()
		if temp > float64(maxEfficiency) && !helpers.StringInSlice(node.Name, exceptionNodesName) {
			maxEfficiency = temp
			mostEfficientNode = node
		}
	}

	if maxEfficiency == 0 {
		return &mostEfficientNode, false
	}

	return &mostEfficientNode, true
}
