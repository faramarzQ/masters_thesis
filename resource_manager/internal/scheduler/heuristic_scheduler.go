package scheduler

import (
	"context"
	"resource_manager/internal/cluster"
	"resource_manager/internal/consts"

	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/klog"
	"k8s.io/kubernetes/pkg/scheduler/framework"
)

type HeuristicScheduler struct {
	*baseScheduler
	framework.FilterPlugin
	framework.PreScorePlugin
	framework.ScorePlugin
}

type HeuristicSchedulerScoreExtension struct{}

func NewHeuristicScheduler() *HeuristicScheduler {
	return &HeuristicScheduler{
		baseScheduler: newBaseScheduler(),
	}
}

func (rs *HeuristicScheduler) Name() string {
	return consts.HEURISTIC_SCHEDULER
}

func (rs *HeuristicScheduler) factory(_ runtime.Object, _ framework.Handle) (framework.Plugin, error) {
	return rs, nil
}

func (rs *HeuristicScheduler) ScoreExtensions() framework.ScoreExtensions {
	return &HeuristicSchedulerScoreExtension{}
}

func (rs *HeuristicScheduler) Filter(ctx context.Context, state *framework.CycleState, p *v1.Pod, nodeInfo *framework.NodeInfo) *framework.Status {
	pod := cluster.BindPod(*p)
	node := cluster.BindNode(*nodeInfo.Node())
	var status framework.Code = framework.Success

	status = rs.baseFilter(pod, node)

	klog.Info("Filtering Pod: ", pod.Name, " on ", nodeInfo.Node().Name, " : ", status.String())
	return framework.NewStatus(status)
}

func (rs *HeuristicScheduler) PreScore(ctx context.Context, state *framework.CycleState, p *v1.Pod, rawNodes []*v1.Node) *framework.Status {
	pod := cluster.BindPod(*p)
	var nodes cluster.NodeList
	for _, rawNode := range rawNodes {
		nodes = append(nodes, cluster.BindNode(*rawNode))
	}
	var status framework.Code = framework.Success

	var nodesScore = make(map[string]int64)
	for _, node := range nodes {
		nodesScore[node.Name] = int64(node.GetCpuEfficiency() * 100)
	}
	rs.set("nodes_score_for_"+pod.Name, nodesScore)

	klog.Info("PreScoring Pod: ", pod.Name)
	return framework.NewStatus(status)
}

func (rs *HeuristicScheduler) Score(ctx context.Context, state *framework.CycleState, p *v1.Pod, nodeName string) (int64, *framework.Status) {
	pod := cluster.BindPod(*p)

	nodesScore := rs.get("nodes_score_for_" + pod.Name).(map[string]int64)
	score := nodesScore[nodeName]

	klog.Info("Scoring Pod: ", pod.Name, " on ", nodeName, " : ", score)
	return score, framework.NewStatus(framework.Success)
}

func (se *HeuristicSchedulerScoreExtension) NormalizeScore(ctx context.Context, state *framework.CycleState, p *v1.Pod, scores framework.NodeScoreList) *framework.Status {
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
