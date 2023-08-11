package scaler

import (
	"bytes"
	"encoding/json"
	"errors"
	"math"
	"net/http"
	"os"
	"resource_manager/internal/cluster"
	"resource_manager/internal/consts"
	"resource_manager/internal/database/repository"

	"k8s.io/klog"
)

type ProposedScaler struct {
	baseScaler
}

type AIServerResponse struct {
	Action  int8
	Epsilon uint8
}

func NewProposedScaler() *ProposedScaler {
	return &ProposedScaler{}
}

func (rs *ProposedScaler) getName() string {
	return consts.PROPOSED_SCALER
}

// TODO: return error
func (ps *ProposedScaler) planScaling(clusterMetrics cluster.ClusterMetrics) error {
	klog.Info(consts.MSG_RUNNING_SCALE_PLANNING)
	defer klog.Info(consts.MSG_FINISHED_SCALE_PLANNING)

	// collect metrics
	// TODO: move to cluster module
	clusterStatus := cluster.GetClusterStatus()
	clusterStatus.Step = scalerExecutionLog.Step
	if previousScalerExecutionLog != nil {
		clusterStatus.PreviousState = (*previousScalerExecutionLog).ScalerExecutionLogDetails.State
		clusterStatus.PreviousAction = (*previousScalerExecutionLog).ScalerExecutionLogDetails.ActionTaken
		clusterStatus.ClusterEnergyConsumption = cluster.CalculateEnergyConsumption(*previousScalerExecutionLog)
		clusterStatus.SuccessRequestRate = cluster.GetSuccessRequestRate()
	}

	payload, err := json.Marshal(clusterStatus)
	if err != nil {
		return err
	}

	response, err := http.Post(os.Getenv("AI_SERVER_URL"), "application/json",
		bytes.NewBuffer(payload))
	if err != nil {
		return err
	}

	if response.StatusCode == http.StatusInternalServerError {
		return errors.New("Received 500 internal server error from the AI-agent!")
	}

	var responseMap map[string]interface{}
	json.NewDecoder(response.Body).Decode(&responseMap)

	repository.InsertScalerExecutionLogDetail(
		scalerExecutionLog,
		string(responseMap["state"].(string)),
		int8(responseMap["action"].(float64)),
	)

	ps.ScaleNodesBetweenOffAndIdleClasses(int8(responseMap["action"].(float64)))

	return nil
}

func (ps *ProposedScaler) ScaleNodesBetweenOffAndIdleClasses(numberOfNodesToScale int8) {
	if numberOfNodesToScale > 0 {
		ps.ScaleNodesFromOffToIdleClass(numberOfNodesToScale)
	} else {
		ps.ScaleNodesFromIdleToOffClass(int8(math.Abs(float64(numberOfNodesToScale))))
	}
}

func (ps *ProposedScaler) ScaleNodesFromOffToIdleClass(numberOfNodesToScale int8) {
	var selectedNodes cluster.NodeList

	for i := 0; i < int(numberOfNodesToScale); i++ {
		if i%2 == 0 {
			efficientNode, ok := cluster.GetMostMemoryEfficientNode(selectedNodes.Names(), consts.OFF_CLASS)
			if !ok {
				break
			}
			selectedNodes = append(selectedNodes, *efficientNode)

		} else {
			efficientNode, ok := cluster.GetMostMemoryEfficientNode(selectedNodes.Names(), consts.OFF_CLASS)
			if !ok {
				break
			}
			selectedNodes = append(selectedNodes, *efficientNode)
		}
	}

	if len(selectedNodes) != 0 {
		var nodeTransition nodeTransition
		nodeTransition.from = consts.OFF_CLASS
		nodeTransition.to = consts.IDLE_CLASS
		nodeTransition.nodesList = selectedNodes
		ps.setTransitions(nodeTransition)
	}
}

func (ps *ProposedScaler) ScaleNodesFromIdleToOffClass(numberOfNodesToScale int8) {
	var selectedNodes cluster.NodeList

	for i := 0; i < int(numberOfNodesToScale); i++ {
		if i%2 == 0 {
			efficientNode, ok := cluster.GetMostMemoryEfficientNode(selectedNodes.Names(), consts.IDLE_CLASS)
			if !ok {
				break
			}
			selectedNodes = append(selectedNodes, *efficientNode)

		} else {
			efficientNode, ok := cluster.GetMostMemoryEfficientNode(selectedNodes.Names(), consts.IDLE_CLASS)
			if !ok {
				break
			}
			selectedNodes = append(selectedNodes, *efficientNode)
		}
	}

	if len(selectedNodes) != 0 {
		var nodeTransition nodeTransition
		nodeTransition.from = consts.IDLE_CLASS
		nodeTransition.to = consts.OFF_CLASS
		nodeTransition.nodesList = selectedNodes
		ps.setTransitions(nodeTransition)
	}
}

func (ps ProposedScaler) onFail(err error) {
	klog.Info(consts.MSG_FAILED_PLANNING)
	klog.Error(err)
	repository.DeleteScalingExecutionLog(&scalerExecutionLog)
}
