package scheduler

import (
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/kubernetes/pkg/scheduler/framework"
)

type schedulerInterface interface {
	Name() string
	factory(_ runtime.Object, _ framework.Handle) (framework.Plugin, error)
}
