package scheduler

import (
	"resource_manager/internal/cluster"
	"resource_manager/internal/consts"
	"time"

	"k8s.io/kubernetes/pkg/scheduler/framework"
)

type baseScheduler struct {
	state map[string]any
}

func newBaseScheduler() *baseScheduler {
	return &baseScheduler{map[string]any{}}
}

func (bs *baseScheduler) set(key string, value any) {
	if bs.state == nil {
		bs.state = map[string]any{}
	}

	bs.state[key] = value
}

func (bs *baseScheduler) get(key string) any {
	return bs.state[key]
}

// A filter functionality which should be applied on all filter plugins
func (bs *baseScheduler) baseFilter(pod cluster.Pod, node cluster.Node) framework.Code {
	if node.Class == consts.IDLE_CLASS {
		scaledAt := node.GetScaledAt()
		if scaledAt.Add(time.Minute*time.Duration(consts.IDLE_NODE_DURATION_MINUTES)).Unix() < time.Now().Unix() {
			return framework.Unschedulable
		}
	}

	return framework.Success
}
