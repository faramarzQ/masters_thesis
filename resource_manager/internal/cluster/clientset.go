package cluster

import (
	"resource_manager/internal/config"
	"resource_manager/internal/consts"
	"sync"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/klog"
	metricsv "k8s.io/metrics/pkg/client/clientset/versioned"
)

var once sync.Once
var Clientset *kubernetes.Clientset
var MetricsClientset *metricsv.Clientset

func RegisterClientSet() {
	once.Do(func() {
		configDir := config.CLUSTER_AUTH_CONFIG

		config, err := clientcmd.BuildConfigFromFlags("", configDir)
		if err != nil {
			klog.Exit(consts.ERROR_CREATING_CLIENTSET, err)
		}

		clientset, err := kubernetes.NewForConfig(config)
		if err != nil {
			klog.Exit(consts.ERROR_CREATING_CLIENTSET, err)
		}

		MetricsClientset = metricsv.NewForConfigOrDie(config)

		Clientset = clientset
		klog.Info(consts.MSG_CLIENTSET_CREATED)
	})
}
