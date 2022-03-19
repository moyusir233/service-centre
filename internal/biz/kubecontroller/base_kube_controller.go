package kubecontroller

import (
	"context"
	"github.com/go-kratos/kratos/v2/errors"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	client_metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	label "k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/watch"
	client_appsv1 "k8s.io/client-go/applyconfigurations/apps/v1"
	client_corev1 "k8s.io/client-go/applyconfigurations/core/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"time"
)

// 创建k8s对象时使用的fieldManger的名称
const fieldManager = "service-center"

type baseKubeController struct {
	client    *kubernetes.Clientset
	namespace string
}

func newBaseKubeController(namespace string) (*baseKubeController, error) {
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

	return &baseKubeController{client: clientset, namespace: namespace}, nil
}

// CreateConfigMap 创建指定的configMap
func (c *baseKubeController) CreateConfigMap(name string, labels map[string]string, data map[string]string) (
	*corev1.ConfigMap, error) {
	// 创建configMap的apply配置
	configMapApplyConfiguration := client_corev1.ConfigMap(name, c.namespace).
		WithLabels(labels).WithData(data)

	return c.client.CoreV1().ConfigMaps(c.namespace).Apply(
		context.Background(),
		configMapApplyConfiguration,
		client_metav1.ApplyOptions{
			FieldManager: fieldManager,
		},
	)
}

// CreateService 创建指定的service
func (c *baseKubeController) CreateService(
	name string, labels map[string]string,
	spec *client_corev1.ServiceSpecApplyConfiguration,
) (*corev1.Service, error) {
	serviceApplyConfiguration := client_corev1.Service(name, c.namespace).
		WithLabels(labels).WithSpec(spec)

	return c.client.CoreV1().Services(c.namespace).Apply(
		context.Background(),
		serviceApplyConfiguration,
		client_metav1.ApplyOptions{
			FieldManager: fieldManager,
		},
	)
}

// CreateDeployment 创建指定的deployment，并执行watch直到deployment能够提供服务
func (c *baseKubeController) CreateDeployment(
	name string, labels map[string]string,
	timeout time.Duration,
	spec *client_appsv1.DeploymentSpecApplyConfiguration) (*appsv1.Deployment, error) {
	deploymentApplyConfiguration := client_appsv1.Deployment(name, c.namespace).
		WithLabels(labels).WithSpec(spec)

	deployment, err := c.client.AppsV1().Deployments(c.namespace).Apply(
		context.Background(),
		deploymentApplyConfiguration,
		client_metav1.ApplyOptions{
			FieldManager: fieldManager,
		},
	)
	if err != nil {
		return nil, err
	}

	// 构建标签选择器和字段选择器，筛选出然后执行watch操作
	// 标签选择器和字段选择器的格式参考:
	// https://kubernetes.io/zh/docs/concepts/overview/working-with-objects/labels/
	// https://kubernetes.io/zh/docs/concepts/overview/working-with-objects/field-selectors/
	// watch响应的object的包含的字段与格式参考:
	// https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.22/#watch-deployment-v1-apps
	formatLabels := label.FormatLabels(labels)
	w, err := c.client.AppsV1().Deployments(c.namespace).Watch(
		context.Background(),
		client_metav1.ListOptions{
			TypeMeta:             deployment.TypeMeta,
			LabelSelector:        formatLabels,
			FieldSelector:        "",
			Watch:                true,
			AllowWatchBookmarks:  false,
			ResourceVersion:      "",
			ResourceVersionMatch: "",
			TimeoutSeconds:       nil,
			Limit:                0,
			Continue:             "",
		},
	)
	if err != nil {
		return nil, err
	}
	defer w.Stop()

	// 设置等待更新的计时器为五分钟
	timer := time.After(timeout)
	for {
		select {
		case <-timer:
			// 超时，删除创建的deployment
			c.DeleteResource(deployment.Name, deployment.TypeMeta)

			return nil, errors.New(
				500, "CREATE_DEPLOYMENT_FAIL", "failed to create the deployment")
		case event := <-w.ResultChan():
			if event.Type == watch.Error || event.Type == watch.Deleted {
				return nil, errors.New(
					500, "CREATE_DEPLOYMENT_FAIL", "failed to create the deployment")
			}

			// 当deployment的所有pod都ready时，表示deployment创建成功
			deployStatus := event.Object.(*appsv1.Deployment).Status
			if deployStatus.ReadyReplicas == *deployment.Spec.Replicas {
				return deployment, nil
			}
		}
	}
}

// CreateStatefulSet 创建指定的statefulSet,并执行watch直到statefulSet能够提供服务
func (c *baseKubeController) CreateStatefulSet(
	name string, labels map[string]string,
	timeout time.Duration,
	spec *client_appsv1.StatefulSetSpecApplyConfiguration) (*appsv1.StatefulSet, error) {
	statefulSetApplyConfiguration := client_appsv1.StatefulSet(name, c.namespace).
		WithLabels(labels).WithSpec(spec)
	statefulSet, err := c.client.AppsV1().StatefulSets(c.namespace).Apply(
		context.Background(),
		statefulSetApplyConfiguration,
		client_metav1.ApplyOptions{
			FieldManager: fieldManager,
		},
	)
	if err != nil {
		return nil, err
	}

	// 构建标签选择器和字段选择器，筛选出然后执行watch操作
	// watch响应的object的包含的字段与格式参考:
	// https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.22/#statefulset-v1-apps
	formatLabels := label.FormatLabels(labels)
	w, err := c.client.AppsV1().StatefulSets(c.namespace).Watch(
		context.Background(),
		client_metav1.ListOptions{
			TypeMeta:             statefulSet.TypeMeta,
			LabelSelector:        formatLabels,
			FieldSelector:        "",
			Watch:                true,
			AllowWatchBookmarks:  false,
			ResourceVersion:      "",
			ResourceVersionMatch: "",
			TimeoutSeconds:       nil,
			Limit:                0,
			Continue:             "",
		},
	)
	if err != nil {
		return nil, err
	}
	defer w.Stop()

	// 设置等待更新的计时器为五分钟
	timer := time.After(timeout)
	for {
		select {
		case <-timer:
			// 超时，删除创建的statefulSet
			c.DeleteResource(statefulSet.Name, statefulSet.TypeMeta)

			return nil, errors.New(
				500, "CREATE_STATEFULSET_FAIL", "failed to create the statefulSet")
		case event := <-w.ResultChan():
			if event.Type == watch.Error || event.Type == watch.Deleted {
				return nil, errors.New(
					500, "CREATE_STATEFULSET_FAIL", "failed to create the statefulSet")
			}

			// 当statefulSet的所有pod都ready时，表示statefulSet创建成功
			status := event.Object.(*appsv1.StatefulSet).Status
			if status.ReadyReplicas == *statefulSet.Spec.Replicas {
				return statefulSet, nil
			}
		}
	}
}

// DeleteResource 删除指定的k8s资源
func (c *baseKubeController) DeleteResource(name string, meta client_metav1.TypeMeta) error {
	switch meta.Kind {
	case "Service":
		return c.client.CoreV1().Services(c.namespace).Delete(
			context.Background(),
			name,
			client_metav1.DeleteOptions{
				TypeMeta: meta,
			},
		)
	case "ConfigMap":
		return c.client.CoreV1().ConfigMaps(c.namespace).Delete(
			context.Background(),
			name,
			client_metav1.DeleteOptions{
				TypeMeta: meta,
			},
		)
	case "Deployment":
		return c.client.AppsV1().Deployments(c.namespace).Delete(
			context.Background(),
			name,
			client_metav1.DeleteOptions{
				TypeMeta: meta,
			},
		)
	case "StatefulSet":
		return c.client.AppsV1().StatefulSets(c.namespace).Delete(
			context.Background(),
			name,
			client_metav1.DeleteOptions{
				TypeMeta: meta,
			},
		)
	default:
		return errors.New(
			500, "UNKNOWN_RESOURCE_TYPE", "can't delete the unknown resource")
	}
}

// DeleteResources 依据label批量删除指定的k8s资源
func (c *baseKubeController) DeleteResources(meta client_metav1.TypeMeta, labelSelector string) error {
	switch meta.Kind {
	case "ConfigMap":
		return c.client.CoreV1().ConfigMaps(c.namespace).DeleteCollection(
			context.Background(),
			client_metav1.DeleteOptions{TypeMeta: meta},
			client_metav1.ListOptions{
				TypeMeta:      meta,
				LabelSelector: labelSelector,
			},
		)
	case "Deployment":
		return c.client.AppsV1().Deployments(c.namespace).DeleteCollection(
			context.Background(),
			client_metav1.DeleteOptions{TypeMeta: meta},
			client_metav1.ListOptions{
				TypeMeta:      meta,
				LabelSelector: labelSelector,
			},
		)
	case "StatefulSet":
		return c.client.AppsV1().StatefulSets(c.namespace).DeleteCollection(
			context.Background(),
			client_metav1.DeleteOptions{TypeMeta: meta},
			client_metav1.ListOptions{
				TypeMeta:      meta,
				LabelSelector: labelSelector,
			},
		)
	default:
		return errors.New(
			500, "UNKNOWN_RESOURCE_TYPE", "can't delete the unknown resource")
	}
}
