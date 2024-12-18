package scaler

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"resource_manager/internal/cluster"
	"resource_manager/internal/consts"

	"k8s.io/klog"
)

type SilencerScaler struct {
	baseScaler
}

func NewSilencerScaler() *SilencerScaler {
	return &SilencerScaler{}
}

func (rs *SilencerScaler) getName() string {
	return consts.SILENCER_SCALER
}

func (rs *SilencerScaler) planScaling(clusterMetrics cluster.ClusterMetrics) error {
	klog.Info(consts.MSG_RUNNING_SCALE_PLANNING)
	defer klog.Info(consts.MSG_FINISHED_SCALE_PLANNING)

	err := rs.silenceActiveNodes()
	if err != nil {
		return err
	}

	err = rs.silenceIdleNodes()
	if err != nil {
		return err
	}

	return nil
}

// For every active node, checks if it has been inactive for a while then silences them to lower level classes
func (sc *SilencerScaler) silenceActiveNodes() error {
	klog.Info("Silencing active nodes")

	nodes := cluster.ListActiveNodes()

	// For every scaler type, directly change nodes to off class,
	// but for the proposed scaler, the idle class is also usable
	activeScaler := cluster.MasterNode().Annotations[consts.ACTIVE_SCALER_LABEL_NAME]
	targetClass := consts.OFF_CLASS
	if activeScaler == consts.PROPOSED_SCALER {
		targetClass = consts.IDLE_CLASS
	}

	allPods := cluster.ListAllPods()

	if len(allPods) == 1 {
		for _, pod := range allPods {
			pod.UnsetWarm()

			klog.Info("Unset warm label from ", pod.Name+" pod on "+pod.Spec.NodeName+" node")
		}
		return nil
	}

	var nodesToTransit cluster.NodeList
	for _, node := range nodes {
		podCount := len(node.ListPods())

		// If node has more than one pod or is utilized, unset warm pods and pass
		if podCount > 1 || node.GetCpuUtilization() > 1 {
			for _, pod := range node.ListPods() {
				if pod.IsAlreadyWarm() {
					klog.Info("Unset warm label from ", pod.Name+" pod on "+node.Name+" node")
					pod.UnsetWarm()
				}
			}

			continue
		}

		// If node has no pod, idle it
		if podCount == 0 {
			nodesToTransit = append(nodesToTransit, node)
			klog.Info("Node:" + node.Name + " has no pod!")
			continue
		}

		// If node only has one pod with no utilization
		if node.GetCpuUtilization() < 1 {
			pod := node.ListPods()[0]
			warmPodDuration, err := strconv.Atoi(os.Getenv("SILENCER_SCALER_WARM_POD_DURATION_MINUTES"))
			if err != nil {
				klog.Info(fmt.Sprintf(consts.ERROR_PARSING_ENV, "SILENCER_SCALER_WARM_POD_DURATION_MINUTES"))
				return err
			}

			if pod.IsAlreadyWarm() {
				// if has been warm for a while
				warmedAt := pod.GetWarmedAt()
				if time.Now().After(warmedAt.Add(time.Minute * time.Duration(warmPodDuration))) {
					nodesToTransit = append(nodesToTransit, node)
					pod.UnsetWarm()
					klog.Info("Unset warm label from ", pod.Name+" pod on "+node.Name+" node")
					// Q: so now we have an idle node with a pod on it?
				}

				continue
			} else {
				klog.Info("Set warm label on ", pod.Name+" pod on "+node.Name+" node")
				pod.WarmUp()
				continue
			}
		}
	}

	var nodeTransition nodeTransition
	nodeTransition.from = consts.ACTIVE_CLASS
	nodeTransition.to = targetClass
	nodeTransition.nodesList = nodesToTransit
	sc.setTransitions(nodeTransition)

	klog.Info("Finished silencing active nodes")

	return nil
}

// For every idle node, checks if it has been idle for a while, then silences them to off class
func (sc *SilencerScaler) silenceIdleNodes() error {
	klog.Info("Silencing Idle nodes")

	// Only available for the proposed scaler
	activeScaler := cluster.MasterNode().Labels[consts.ACTIVE_SCALER_LABEL_NAME]
	if activeScaler != consts.PROPOSED_SCALER {
		return nil
	}

	nodes := cluster.ListNodes().InClass(consts.IDLE_CLASS)
	idleNodeDuration, err := strconv.Atoi(os.Getenv("SILENCER_SCALER_IDLE_NODE_DURATION_MINUTES"))
	if err != nil {
		klog.Info(fmt.Sprintf(consts.ERROR_PARSING_ENV, "SILENCER_SCALER_IDLE_NODE_DURATION_MINUTES"))
		return err
	}

	var nodesToTransit cluster.NodeList
	for _, node := range nodes {
		// if node has been idle for a specific time, scale to off
		scaledAt := node.GetScaledAt()
		if scaledAt.Add(time.Minute*time.Duration(idleNodeDuration)).Unix() < time.Now().Unix() {
			nodesToTransit = append(nodesToTransit, node)
		}
	}

	var nodeTransition nodeTransition
	nodeTransition.from = consts.IDLE_CLASS
	nodeTransition.to = consts.OFF_CLASS
	nodeTransition.nodesList = nodesToTransit
	sc.setTransitions(nodeTransition)

	klog.Info("Finished silencing active nodes")

	return nil
}
