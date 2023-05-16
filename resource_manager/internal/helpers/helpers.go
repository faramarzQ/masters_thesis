package helpers

import (
	"math/rand"
	"resource_manager/internal/cluster"
	"time"
)

func GetRandomNodesFromNodeList(nodeList cluster.NodeList, numberOfNodesToSelect int64) cluster.NodeList {
	nodesIndexesToSelect := []int{}

	for len(nodesIndexesToSelect) < int(numberOfNodesToSelect) {
		rand.Seed(time.Now().UnixNano())
		randomNum := rand.Intn(int(len(nodeList)))

		var nodeAlreadySelected bool
		for i := 0; i < len(nodesIndexesToSelect); i++ {
			if randomNum == nodesIndexesToSelect[i] {
				nodeAlreadySelected = true
				break
			}
		}

		if !nodeAlreadySelected {
			nodesIndexesToSelect = append(nodesIndexesToSelect, randomNum)
		}
	}

	var nodesToSelect cluster.NodeList
	for i := 0; i < len(nodesIndexesToSelect); i++ {
		nodesToSelect = append(nodesToSelect, nodeList[nodesIndexesToSelect[i]])
	}

	return nodesToSelect
}
