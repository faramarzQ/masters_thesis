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
	framework.PostBindPlugin
}

type RandomSchedulerScoreExtension struct{}

func NewRandomScheduler() *RandomScheduler {
	return &RandomScheduler{
		baseScheduler: newBaseScheduler(),
	}
}

func (rs *RandomScheduler) Name() string {
	return consts.RANDOM_SCHEDULER
}

func (rs *RandomScheduler) factory(_ runtime.Object, _ framework.Handle) (framework.Plugin, error) {
	return rs, nil
}

func (rs *RandomScheduler) ScoreExtensions() framework.ScoreExtensions {
	return &RandomSchedulerScoreExtension{}
}

func (rs *RandomScheduler) Filter(ctx context.Context, state *framework.CycleState, p *v1.Pod, nodeInfo *framework.NodeInfo) *framework.Status {
	pod := cluster.BindPod(*p)
	node := cluster.BindNode(*nodeInfo.Node())
	klog.Infof("Started Filtering Pod: %s on %s", pod.Name, nodeInfo.Node().Name)

	var status framework.Code = framework.Success

	status = rs.baseFilter(pod, node)

	klog.Infof("Finished Filtering Pod: %s on %s : %s", pod.Name, nodeInfo.Node().Name, status.String())
	return framework.NewStatus(status)
}

func (rs *RandomScheduler) PreScore(ctx context.Context, state *framework.CycleState, p *v1.Pod, rawNodes []*v1.Node) *framework.Status {
	pod := cluster.BindPod(*p)
	klog.Infof("Started PreScoring Pod: %s", pod.Name)

	var nodes cluster.NodeList
	for _, rawNode := range rawNodes {
		nodes = append(nodes, cluster.BindNode(*rawNode))
	}
	var status framework.Code = framework.Success

	// Give random score to every node
	rand.Seed(time.Now().UnixNano())
	var nodesScore = make(map[string]int64)
	for _, node := range nodes {
		nodesScore[node.Name] = int64(rand.Intn(int(framework.MaxNodeScore)))
	}
	rs.set("nodes_score_for_"+pod.Name, nodesScore)

	klog.Infof("Finished PreScoring Pod: %s", pod.Name)

	return framework.NewStatus(status)
}

func (rs *RandomScheduler) Score(ctx context.Context, state *framework.CycleState, p *v1.Pod, nodeName string) (int64, *framework.Status) {
	pod := cluster.BindPod(*p)
	klog.Infof("Started Scoring Pod: %s on %s", pod.Name, nodeName)

	nodesScore := rs.get("nodes_score_for_" + pod.Name).(map[string]int64)
	score := nodesScore[nodeName]

	klog.Infof("Finished Scoring Pod: %s on %s : %d", pod.Name, nodeName, score)
	return score, framework.NewStatus(framework.Success)
}

func (*RandomSchedulerScoreExtension) NormalizeScore(ctx context.Context, state *framework.CycleState, p *v1.Pod, scores framework.NodeScoreList) *framework.Status {
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

func (rs *RandomScheduler) PostBind(ctx context.Context, state *framework.CycleState, p *v1.Pod, nodeName string) {
	klog.Info("Started PostBind.")

	node := cluster.GetNodeByName(nodeName)
	if node.Class != consts.ACTIVE_CLASS {
		node.SetClass(consts.ACTIVE_CLASS)
	}

	klog.Info("Finished PostBind")
}

func max(a, b int64) int64 {
	if a < b {
		return b
	}
	return a
}
