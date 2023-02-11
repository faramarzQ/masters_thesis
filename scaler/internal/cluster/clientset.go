package cluster

import (
	"fmt"
	"log"
	"scaler/config"
	"scaler/internal/consts"
	"sync"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

var once sync.Once
var Clientset *kubernetes.Clientset

func RegisterClientSet() {
	once.Do(func() {
		configDir := config.CLUSTER_AUTH_CONFIG

		config, err := clientcmd.BuildConfigFromFlags("", configDir)
		if err != nil {
			log.Fatal(consts.ERROR_CREATING_CLIENTSET, err)
		}

		clientset, err := kubernetes.NewForConfig(config)
		if err != nil {
			log.Fatal(consts.ERROR_CREATING_CLIENTSET, err)
		}

		Clientset = clientset
		fmt.Println(consts.CLIENTSET_CREATED)
	})
}
