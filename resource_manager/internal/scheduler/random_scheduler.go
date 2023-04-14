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

type RandomScheduler struct {
	framework.FilterPlugin
}

func NewRandomScheduler() *RandomScheduler {
	return &RandomScheduler{}
}

func (rs *RandomScheduler) Name() string {
	return consts.RANDOM_SCHEDULER_NAME
}

func (rs *RandomScheduler) factory(_ runtime.Object, _ framework.Handle) (framework.Plugin, error) {
	return rs, nil
}

func (rs *RandomScheduler) Filter(ctx context.Context, state *framework.CycleState, p *v1.Pod, nodeInfo *framework.NodeInfo) *framework.Status {
	pod := cluster.BindPod(*p)
	// node := cluster.BindNode(*nodeInfo.Node())
	var status framework.Code = framework.Success

	klog.Info("Filtering Pod: ", pod.Name, " on ", nodeInfo.Node().Name, " : ", status.String())
	return framework.NewStatus(status)
}
