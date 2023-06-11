package scaler

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"resource_manager/internal/cluster"
	"resource_manager/internal/consts"

	"k8s.io/klog"
)

type ProposedScaler struct {
	baseScaler
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
	fmt.Println(clusterStatus)
	json_data, err := json.Marshal(clusterStatus)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(string(json_data))

	val := cluster.ClusterStatus{}
	json.Unmarshal(json_data, &val)

	fmt.Println(val)

	resp, err := http.Post("http://localhost:8080", "application/json",
		bytes.NewBuffer(json_data))

	if err != nil {
		log.Fatal(err)
	}

	defer resp.Body.Close()
	var res map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&res)
	fmt.Println(res["json"])
}
