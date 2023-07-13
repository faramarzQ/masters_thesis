package scaler

import (
	"bytes"
	"encoding/json"
	"log"
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

func (ps *ProposedScaler) planScaling(clusterMetrics cluster.ClusterMetrics) {
	klog.Info(consts.MSG_RUNNING_SCALE_PLANNING)
	defer klog.Info(consts.MSG_FINISHED_SCALE_PLANNING)

	// collect metrics
	clusterStatus := cluster.GetClusterStatus()
	if previousScalerExecutionLog != nil {
		clusterStatus.ExecutedPreviously = true
		clusterStatus.PreviousState = (*previousScalerExecutionLog).ScalerExecutionLogDetails.State
		clusterStatus.PreviousActionTaken = (*previousScalerExecutionLog).ScalerExecutionLogDetails.ActionTaken
		clusterStatus.PreviousEpsilonValue = (*previousScalerExecutionLog).ScalerExecutionLogDetails.EpsilonValue
	}

	payload, err := json.Marshal(clusterStatus)
	if err != nil {
		log.Fatal(err)
	}

	response, err := http.Post(os.Getenv("AI_SERVER_URL"), "application/json",
		bytes.NewBuffer(payload))
	if err != nil {
		log.Fatal(err)
	}

	var responseMap map[string]interface{}
	json.NewDecoder(response.Body).Decode(&responseMap)

	repository.InsertScalerExecutionLogDetail(
		scalerExecutionLog,
		string(responseMap["state"].(string)),
		int8(responseMap["action"].(float64)),
		uint8(responseMap["epsilon"].(float64)),
	)

	// take action
}
