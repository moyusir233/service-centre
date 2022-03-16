// Package kubecontroller 负责与k8s api-server交互，创建集群组件
package kubecontroller

import (
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

type KubeController struct {
	client *kubernetes.Clientset
}

func NewKubeController() (*KubeController, error) {
	// 通过service account获得访问api server的config实例
	config, err := rest.InClusterConfig()
	if err != nil {
		return nil, err
	}

	// 建立访问k8s的客户端
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, err
	}

	return &KubeController{client: clientset}, nil
}
