package scheduler

import (
	"context"
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
	return &RandomScheduler{
		baseScheduler: newBaseScheduler(),
	}
}

func (rs *RandomScheduler) Name() string {
	return consts.RANDOM_SCHEDULER_NAME
}

func (rs *RandomScheduler) factory(_ runtime.Object, _ framework.Handle) (framework.Plugin, error) {
	return rs, nil
}

func (rs *RandomScheduler) ScoreExtensions() framework.ScoreExtensions {
	return &ScoreExtension{}
}

func (rs *RandomScheduler) Filter(ctx context.Context, state *framework.CycleState, p *v1.Pod, nodeInfo *framework.NodeInfo) *framework.Status {
	pod := cluster.BindPod(*p)
	node := cluster.BindNode(*nodeInfo.Node())
	var status framework.Code = framework.Success

	status = rs.baseFilter(pod, node)

	klog.Info("Filtering Pod: ", pod.Name, " on ", nodeInfo.Node().Name, " : ", status.String())
	return framework.NewStatus(status)
}

func (rs *RandomScheduler) PreScore(ctx context.Context, state *framework.CycleState, p *v1.Pod, rawNodes []*v1.Node) *framework.Status {
	pod := cluster.BindPod(*p)
	var nodes cluster.NodeList
	for _, rawNode := range rawNodes {
		nodes = append(nodes, cluster.BindNode(*rawNode))
	}
	var status framework.Code = framework.Success

	// Give random score to every node
	rand.Seed(time.Now().UnixNano())
	rs.set("selected_nodes_for_"+pod.Name, []string{
		nodes[rand.Intn(len(nodes))].Name,
	})

	var nodesScore = make(map[string]int64)
	for _, node := range nodes {
		nodesScore[node.Name] = int64(rand.Intn(int(framework.MaxNodeScore)))
	}

	klog.Info("PreScoring Pod: ", pod.Name)
	return framework.NewStatus(status)
}

func (rs *RandomScheduler) Score(ctx context.Context, state *framework.CycleState, p *v1.Pod, nodeName string) (int64, *framework.Status) {
	pod := cluster.BindPod(*p)
	var score int64
	nodesScore := rs.get("nodes_score_for_" + pod.Name).(map[string]int64)
	for key, value := range nodesScore {
		if key == nodeName {
			score = value
		}
	}

	klog.Info("Scoring Pod: ", pod.Name, " on ", nodeName, " : ", score)
	return score, framework.NewStatus(framework.Success)
}

func (se *ScoreExtension) NormalizeScore(ctx context.Context, state *framework.CycleState, p *v1.Pod, scores framework.NodeScoreList) *framework.Status {
	var highest int64
	for _, nodeScore := range scores {
		highest = max(int64(highest), nodeScore.Score)
	}

	for i, nodeScore := range scores {
		scores[i].Score = nodeScore.Score * int64(framework.MaxNodeScore) / highest
	}

	klog.Info("Normalized scores: ", scores)
	return framework.NewStatus(framework.Success)
}

func max(a, b int64) int64 {
	if a < b {
		return b
	}
	return a
}
