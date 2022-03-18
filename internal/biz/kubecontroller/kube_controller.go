// Package kubecontroller 负责与k8s api-server交互，创建集群组件
package kubecontroller

import (
	"fmt"
	corev1 "k8s.io/api/core/v1"
)

type KubeController struct {
	*baseKubeController
}

func NewKubeController(namespace string) (*KubeController, error) {
	controller, err := newBaseKubeController(namespace)
	if err != nil {
		return nil, err
	}

	return &KubeController{baseKubeController: controller}, nil
}

// CreateConfigMapOfGeneratedCode 为dataCollection与dataProcessing服务生成的代码创建configMap
func (c *KubeController) CreateConfigMapOfGeneratedCode(username string, dcCode, dpCode map[string]string) (
	dcCm *corev1.ConfigMap, dpCm *corev1.ConfigMap, err error) {
	// 保存生成代码的cm统一以<用户名>-<服务名简写,dc或dp>-code命名
	// 并以user:username作为label
	dcName := fmt.Sprintf("%s-%s-code", username, "dc")
	dpName := fmt.Sprintf("%s-%s-code", username, "dp")
	label := map[string]string{"user": username}

	dcCm, err = c.CreateConfigMap(dcName, label, dcCode)
	if err != nil {
		return nil, nil, err
	}
	dpCm, err = c.CreateConfigMap(dpName, label, dpCode)
	if err != nil {
		return nil, nil, err
	}
	return
}

// CreateConfigMapOfStateRegisterInfo 创建保存设备状态注册信息的configMap
func (c *KubeController) CreateConfigMapOfStateRegisterInfo() {

}
