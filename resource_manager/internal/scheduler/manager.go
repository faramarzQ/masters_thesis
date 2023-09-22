package scheduler

import (
	"os"
	"resource_manager/internal/cluster"
	"resource_manager/internal/config"
	"resource_manager/internal/consts"

	"github.com/spf13/cobra"
	"k8s.io/klog"
	"k8s.io/kubernetes/cmd/kube-scheduler/app"

	"k8s.io/component-base/logs"
)

type schedulerManager struct {
	schedulerCommand *cobra.Command
	scheduler        schedulerInterface
}

func NewSchedulerManager() *schedulerManager {
	schedulerManager := schedulerManager{}

	schedulerManager.RegisterActiveScheduler(
		NewRandomScheduler(),
		NewHeuristicScheduler(),
	)

	command := app.NewSchedulerCommand(
		app.WithPlugin(schedulerManager.scheduler.Name(), schedulerManager.scheduler.factory),
	)
	schedulerManager.schedulerCommand = command

	logs.InitLogs()
	defer logs.FlushLogs()

	return &schedulerManager
}

// Registers the active scheduler into the scheduler manager
func (sm *schedulerManager) RegisterActiveScheduler(schedulers ...schedulerInterface) {
	for _, scheduler := range schedulers {
		if scheduler.Name() == config.ACTIVE_SCHEDULER {
			sm.scheduler = scheduler
			break
		}
	}
	if sm.scheduler == nil {
		klog.Error(consts.ERROR_REGISTERING_SCHEDULER)
	}

	klog.Infof(consts.MSG_REGISTERED_ACTIVE_SCHEDULER, sm.scheduler.Name())
}

func (sm *schedulerManager) Run() {
	cluster.LabelNewNodes()

	if err := sm.schedulerCommand.Execute(); err != nil {
		os.Exit(1)
	}
}
