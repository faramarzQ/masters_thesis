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

	status = rs.baseFilter(pod, node)

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
	rs.set("selected_nodes_for_"+pod.Name, []string{
		nodes[rand.Intn(len(nodes))].Name,
		nodes[rand.Intn(len(nodes))].Name,
		nodes[rand.Intn(len(nodes))].Name,
	})

	klog.Info("PreScoring Pod: ", pod.Name)
	return framework.NewStatus(status)
}

func (rs *RandomScheduler) Score(ctx context.Context, state *framework.CycleState, p *v1.Pod, nodeName string) (int64, *framework.Status) {
	pod := cluster.BindPod(*p)
	// var status framework.Code = framework.Success
	var score int64

	mostPrioritizedNodes := (rs.get("selected_nodes_for_" + pod.Name)).([]string)
	fmt.Println(mostPrioritizedNodes)
	for i := 0; i < len(mostPrioritizedNodes); i++ {
		if mostPrioritizedNodes[i] == nodeName {
			score = 20
		}
	}

	klog.Info("Scoring Pod: ", pod.Name, " on ", nodeName, " : ", score)
	return score, nil
}

func (rs *ScoreExtension) ScoreExtensions() framework.ScoreExtensions {
	return rs
}

func (se *ScoreExtension) NormalizeScore(ctx context.Context, state *framework.CycleState, p *v1.Pod, scores framework.NodeScoreList) *framework.Status {
	var highest int64
	for _, nodeScore := range scores {
		highest = max(int64(highest), nodeScore.Score)
	}
	MaxNodeScore := 20
	for i, nodeScore := range scores {
		scores[i].Score = nodeScore.Score * int64(MaxNodeScore) / highest
	}

	fmt.Println("Normalized scores: ", scores)

	return nil
}

func max(a, b int64) int64 {
	if a < b {
		return b
	}
	return a
}
