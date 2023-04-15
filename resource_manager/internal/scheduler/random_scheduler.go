package scheduler

import (
	"context"
	"fmt"
	"math/rand"
	"resource_manager/internal/cluster"
	"resource_manager/internal/consts"
	"time"

	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/klog"
	"k8s.io/kubernetes/pkg/scheduler/framework"
)

type RandomScheduler struct {
	*baseScheduler
	framework.FilterPlugin
	framework.PreScorePlugin
	framework.ScorePlugin
}

type ScoreExtension struct{}

func NewRandomScheduler() *RandomScheduler {
	return &RandomScheduler{baseScheduler: newBaseScheduler()}
}

func (rs *RandomScheduler) Name() string {
	return consts.RANDOM_SCHEDULER_NAME
}

func (rs *RandomScheduler) factory(_ runtime.Object, _ framework.Handle) (framework.Plugin, error) {
	return rs, nil
}

func (rs *RandomScheduler) Filter(ctx context.Context, state *framework.CycleState, p *v1.Pod, nodeInfo *framework.NodeInfo) *framework.Status {
	pod := cluster.BindPod(*p)
	node := cluster.BindNode(*nodeInfo.Node())
	var status framework.Code = framework.Success

	// Only active nodes are available
	if node.Class != consts.ACTIVE_CLASS {
		status = framework.Unschedulable
	}

	klog.Info("Filtering Pod: ", pod.Name, " on ", nodeInfo.Node().Name, " : ", status.String())
	return framework.NewStatus(status)
}

func (rs *RandomScheduler) PreScore(ctx context.Context, state *framework.CycleState, p *v1.Pod, rawNodes []*v1.Node) *framework.Status {
	pod := cluster.BindPod(*p)
	var nodes cluster.NodeList
	for idx, rawNode := range rawNodes {
		fmt.Println("idx", idx, rawNode.Name)
		nodes = append(nodes, cluster.BindNode(*rawNode))
	}
	var status framework.Code = framework.Success

	// Select a random node
	rand.Seed(time.Now().UnixNano())
	randomIndex := rand.Intn(len(nodes))
	fmt.Println("indx", randomIndex)
	rs.set("selected_node_for_"+pod.Name, nodes[randomIndex].Name)
	fmt.Println(nodes[randomIndex].Name)

	klog.Info("PreScoring Pod: ", pod.Name)
	return framework.NewStatus(status)
}

func (rs *RandomScheduler) Score(ctx context.Context, state *framework.CycleState, p *v1.Pod, nodeName string) (int64, *framework.Status) {
	pod := cluster.BindPod(*p)
	var status framework.Code = framework.Success
	var score int64

	if rs.get("selected_node_for_"+pod.Name) == nodeName {
		score = 1
	}

	klog.Info("Scoring Pod: ", pod.Name, " on ", nodeName, " : ", score)
	return score, framework.NewStatus(status)
}

func (rs *RandomScheduler) ScoreExtensions() framework.ScoreExtensions {
	return ScoreExtension{}
}

func (se ScoreExtension) NormalizeScore(ctx context.Context, state *framework.CycleState, p *v1.Pod, scores framework.NodeScoreList) *framework.Status {
	// pod := cluster.BindPod(*p)
	var status framework.Code = framework.Success

	fmt.Println(scores)
	// klog.Info("Scoring Pod: ", pod.Name, " : ", status.String())
	return framework.NewStatus(status)
}
