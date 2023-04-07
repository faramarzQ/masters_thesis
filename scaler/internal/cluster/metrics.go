package cluster

type ClusterMetrics struct {
	ActiveNodesMetrics ActiveNodesMetrics
}

func GetClusterMetrics() ClusterMetrics {
	return ClusterMetrics{
		GetActiveNodesMetrics(),
	}
}

// Metrics of cluster nodes
type NodeMetrics struct {
	CpuUtilization    float64
	MemoryUtilization float64
}

// Calculated metric values for on every active nodes
type ActiveNodesMetrics map[string]NodeMetrics

func GetActiveNodesMetrics() ActiveNodesMetrics {
	activeNodes := ListActiveNodes()
	var activeNodesMetrics ActiveNodesMetrics = make(ActiveNodesMetrics, len(activeNodes))
	for _, node := range activeNodes {
		cpuUtilization := node.GetCpuUtilization()
		memoryUtilization := node.GetMemoryUtilization()

		activeNodesMetrics[node.Name] = NodeMetrics{
			CpuUtilization:    cpuUtilization,
			MemoryUtilization: memoryUtilization,
		}
	}
	return activeNodesMetrics
}

func (cm ClusterMetrics) GetAverageCpuUtilization() float64 {
	var sumCpuUtil float64
	for _, resourceMetrics := range cm.ActiveNodesMetrics {
		sumCpuUtil += resourceMetrics.CpuUtilization
	}
	return sumCpuUtil / float64(len(cm.ActiveNodesMetrics))
}

func (cm ClusterMetrics) GetAverageMemoryUtilization() float64 {
	var sumMemoryUtil float64
	for _, resourceMetrics := range cm.ActiveNodesMetrics {
		sumMemoryUtil += resourceMetrics.MemoryUtilization
	}
	return sumMemoryUtil / float64(len(cm.ActiveNodesMetrics))
}
