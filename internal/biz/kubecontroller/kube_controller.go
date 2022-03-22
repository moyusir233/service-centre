// Package kubecontroller 负责与k8s api-server交互，创建集群组件
package kubecontroller

import (
	"fmt"
	v1 "gitee.com/moyusir/util/api/util/v1"
	"github.com/go-kratos/kratos/v2/errors"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
	"k8s.io/apimachinery/pkg/util/json"
	client_appsv1 "k8s.io/client-go/applyconfigurations/apps/v1"
	client_corev1 "k8s.io/client-go/applyconfigurations/core/v1"
	client_metav1 "k8s.io/client-go/applyconfigurations/meta/v1"
	"k8s.io/utils/pointer"
	"time"
)

type KubeController struct {
	*baseKubeController
}

// BaseDeployOption 部署时的基本配置
type BaseDeployOption struct {
	// 用户名
	Username string
	// 启动的副本数量
	Replica int32
	// 超时时长
	Timeout time.Duration

	// 服务运行时配置
	// 运行服务使用的镜像
	Image string
	// 运行服务的可执行程序对应的cm
	Exe *corev1.ConfigMap
	// 对应服务所需的可执行程序的数据在cm data中的key
	ExeItemKey string
}
type DataProcessingDeployOption struct {
	BaseDeployOption
	// 设备状态注册信息对应的cm
	RegisterInfo *corev1.ConfigMap
}

type DataCollectionDeployOption struct {
	BaseDeployOption
	// 项目部署时使用的域名，用于路由匹配
	AppDomainName string
}

func NewKubeController(namespace string) (*KubeController, error) {
	controller, err := newBaseKubeController(namespace)
	if err != nil {
		return nil, err
	}

	return &KubeController{baseKubeController: controller}, nil
}

// Unregister 清空用户相关的k8s资源
func (c *KubeController) Unregister(username string) error {
	labelSelector := "user=" + username
	types := []string{"Deployment", "StatefulSet", "Service", "ConfigMap"}

	for _, resourceType := range types {
		err := c.DeleteResources(resourceType, labelSelector)
		if err != nil {
			return err
		}
	}

	return nil
}

// CreateConfigMapOfExe 创建可执行程序对应的cm，将二进制数据保存进cm中
func (c *KubeController) CreateConfigMapOfExe(username string, exe map[string][]byte) (
	*corev1.ConfigMap, error) {
	// 以user:username为label,以username-exe为名称创建cm
	label := map[string]string{"user": username}
	return c.CreateConfigMap(username+"-exe", label, nil, exe)
}

func (c *KubeController) CreateConfigMapOfStateRegisterInfo(
	username string, info []*v1.DeviceStateRegisterInfo) (*corev1.ConfigMap, error) {
	// 以user:username为label,username-state-register-info为名称创建cm
	// 保存注册信息的json数据
	marshal, err := json.Marshal(info)
	if err != nil {
		return nil, err
	}

	label := map[string]string{"user": username}
	return c.CreateConfigMap(username+"-state-register-info", label, map[string]string{
		"register_info.json": string(marshal),
	}, nil)
}

// DeployDataProcessingService 部署数据处理服务，返回指向应用容器endpoint的service组件的信息，提供给网关注册使用
func (c *KubeController) DeployDataProcessingService(option *DataProcessingDeployOption) (*corev1.Service, error) {
	if option == nil {
		return nil, errors.New(500, "option is nil", "")
	}

	// deployment以<用户名>-dp命名，以app:<用户名>-dp和user:<username>为label
	name := fmt.Sprintf("%s-dp", option.Username)
	label := map[string]string{"app": name, "user": option.Username}

	deploymentSpec := getDataProcessingDeploymentSpec(name, label, option)
	_, err := c.CreateDeployment(name, label, option.Timeout, deploymentSpec)
	if err != nil {
		return nil, err
	}

	// 为deployment创建负责负载均衡的service
	serviceLabel := map[string]string{"user": option.Username}
	serviceType := corev1.ServiceTypeClusterIP
	serviceSpec := client_corev1.ServiceSpecApplyConfiguration{
		Ports: []client_corev1.ServicePortApplyConfiguration{
			{
				Name: pointer.String("http"),
				Port: pointer.Int32(8000),
			},
		},
		Selector: label,
		Type:     &serviceType,
	}

	// TODO service创建失败后，是否要撤销之前deployment的部署？
	return c.CreateService(name, serviceLabel, &serviceSpec)
}

// DeployDataCollectionService 部署数据收集服务,返回指向应用容器endpoint的service组件的信息，提供给网关注册使用
func (c *KubeController) DeployDataCollectionService(option *DataCollectionDeployOption) (*corev1.Service, error) {
	if option == nil {
		return nil, errors.New(500, "option is nil", "")
	}

	// statefulSet以<用户名>-dc命名，以app:<<用户名>-dc>,user:<username>为label
	name := fmt.Sprintf("%s-dc", option.Username)
	label := map[string]string{"app": name, "user": option.Username}

	// 先创建statefulSet所需的无头服务，以<用户名>-dc-headless命名，以user:<username>为label
	headlessServiceName := fmt.Sprintf("%s-dc-headless", option.Username)
	serviceLabel := map[string]string{"user": option.Username}
	serviceType := corev1.ServiceTypeClusterIP
	serviceSpec := client_corev1.ServiceSpecApplyConfiguration{
		Ports: []client_corev1.ServicePortApplyConfiguration{
			{
				Name: pointer.String("http"),
				Port: pointer.Int32(8000),
			},
			{
				Name: pointer.String("grpc"),
				Port: pointer.Int32(9000),
			},
		},
		Selector: label,
		// 创建无头服务所需的选项
		ClusterIP: pointer.String("None"),
		Type:      &serviceType,
	}

	headlessService, err := c.CreateService(headlessServiceName, serviceLabel, &serviceSpec)
	if err != nil {
		return nil, err
	}

	// 创建statefulSet
	statefulSetSpec := getDataCollectionStatefulSetSpec(name, headlessService.Name, label, option)
	_, err = c.CreateStatefulSet(name, label, option.Timeout, statefulSetSpec)
	if err != nil {
		return nil, err
	}

	// 为statefulSet创建负责负载均衡的service，重用headless service的配置
	serviceSpec.ClusterIP = nil
	return c.CreateService(name, serviceLabel, &serviceSpec)
}

// 辅助函数，创建dataProcessing服务的部署配置
func getDataProcessingDeploymentSpec(name string, label map[string]string, option *DataProcessingDeployOption) *client_appsv1.DeploymentSpecApplyConfiguration {
	// 配置部署选项
	var (
		imagePullPolicy = corev1.PullIfNotPresent
		restartPolicy   = corev1.RestartPolicyAlways
		servicePort     = intstr.FromInt(8000)
	)

	return &client_appsv1.DeploymentSpecApplyConfiguration{
		// 配置部署的副本数量和selector
		Replicas: &option.Replica,
		Selector: &client_metav1.LabelSelectorApplyConfiguration{
			MatchLabels: label,
		},

		// 配置pod模板
		Template: &client_corev1.PodTemplateSpecApplyConfiguration{
			// pod使用和其控制器一样的name和label
			ObjectMetaApplyConfiguration: &client_metav1.ObjectMetaApplyConfiguration{
				Name:   &name,
				Labels: label,
			},
			Spec: &client_corev1.PodSpecApplyConfiguration{
				// 配置需要挂载的volume
				Volumes: []client_corev1.VolumeApplyConfiguration{
					{
						// 存放编译得到的二进制程序
						Name: pointer.String("app"),
						VolumeSourceApplyConfiguration: client_corev1.VolumeSourceApplyConfiguration{
							ConfigMap: &client_corev1.ConfigMapVolumeSourceApplyConfiguration{
								LocalObjectReferenceApplyConfiguration: client_corev1.LocalObjectReferenceApplyConfiguration{
									Name: pointer.String(option.Exe.Name),
								},
								Items: []client_corev1.KeyToPathApplyConfiguration{
									{
										Key:  pointer.String(option.ExeItemKey),
										Path: pointer.String("server"),
									},
								},
								DefaultMode: pointer.Int32(0777),
							},
						},
					},
					{
						// dataProcessing服务运行所需的配置文件volume，包括通用配置和用户设备状态注册信息
						Name: pointer.String("config"),
						VolumeSourceApplyConfiguration: client_corev1.VolumeSourceApplyConfiguration{
							// 合并通用配置和用户注册信息的configMap到一个目录下进行挂载
							Projected: &client_corev1.ProjectedVolumeSourceApplyConfiguration{
								Sources: []client_corev1.VolumeProjectionApplyConfiguration{
									{
										// 通用配置
										ConfigMap: &client_corev1.ConfigMapProjectionApplyConfiguration{
											LocalObjectReferenceApplyConfiguration: client_corev1.LocalObjectReferenceApplyConfiguration{
												Name: pointer.String("config"),
											},
											Items: []client_corev1.KeyToPathApplyConfiguration{
												{
													Key:  pointer.String("data-processing-config.yaml"),
													Path: pointer.String("config.yaml"),
												},
											},
										},
									},
									{
										// 用户设备状态注册信息
										ConfigMap: &client_corev1.ConfigMapProjectionApplyConfiguration{
											LocalObjectReferenceApplyConfiguration: client_corev1.LocalObjectReferenceApplyConfiguration{
												Name: &option.RegisterInfo.Name,
											},
										},
									},
								},
							},
						},
					},
				},
				Containers: []client_corev1.ContainerApplyConfiguration{
					{
						Name:            &name,
						Image:           &option.Image,
						ImagePullPolicy: &imagePullPolicy,
						VolumeMounts: []client_corev1.VolumeMountApplyConfiguration{
							// 挂载放置编译后得到的二进制可执行程序的共享目录
							{
								Name:      pointer.String("app"),
								MountPath: pointer.String("/app"),
							},
							// 挂载应用运行所需的配置文件目录
							{
								Name:      pointer.String("config"),
								MountPath: pointer.String("/etc/app-configs"),
							},
						},
						Env: []client_corev1.EnvVarApplyConfiguration{
							{
								Name:  pointer.String("USERNAME"),
								Value: pointer.String(option.Username),
							},
						},
						Ports: []client_corev1.ContainerPortApplyConfiguration{
							{
								ContainerPort: pointer.Int32(8000),
							},
						},
						LivenessProbe: &client_corev1.ProbeApplyConfiguration{
							HandlerApplyConfiguration: client_corev1.HandlerApplyConfiguration{
								TCPSocket: &client_corev1.TCPSocketActionApplyConfiguration{
									Port: &servicePort,
								},
							},
							InitialDelaySeconds: pointer.Int32(20),
							PeriodSeconds:       pointer.Int32(60),
						},
						ReadinessProbe: &client_corev1.ProbeApplyConfiguration{
							HandlerApplyConfiguration: client_corev1.HandlerApplyConfiguration{
								TCPSocket: &client_corev1.TCPSocketActionApplyConfiguration{
									Port: &servicePort,
								},
							},
							InitialDelaySeconds: pointer.Int32(20),
							PeriodSeconds:       pointer.Int32(20),
						},
					},
				},
				RestartPolicy: &restartPolicy,
			},
		},
	}
}

// 辅助函数，创建dataCollection服务的部署配置
func getDataCollectionStatefulSetSpec(
	name, headlessServiceName string,
	label map[string]string,
	option *DataCollectionDeployOption) *client_appsv1.StatefulSetSpecApplyConfiguration {
	// 配置部署选项
	var (
		imagePullPolicy = corev1.PullIfNotPresent
		restartPolicy   = corev1.RestartPolicyAlways
		servicePort     = intstr.FromInt(8000)
	)

	return &client_appsv1.StatefulSetSpecApplyConfiguration{
		// 配置部署的副本数量和selector
		Replicas: &option.Replica,
		Selector: &client_metav1.LabelSelectorApplyConfiguration{
			MatchLabels: label,
		},

		// 配置pod模板
		Template: &client_corev1.PodTemplateSpecApplyConfiguration{
			// pod使用和其控制器一样的name和label
			ObjectMetaApplyConfiguration: &client_metav1.ObjectMetaApplyConfiguration{
				Name:   &name,
				Labels: label,
			},
			Spec: &client_corev1.PodSpecApplyConfiguration{
				// 配置需要挂载的volume
				Volumes: []client_corev1.VolumeApplyConfiguration{
					{
						// 存放编译得到的二进制程序
						Name: pointer.String("app"),
						VolumeSourceApplyConfiguration: client_corev1.VolumeSourceApplyConfiguration{
							ConfigMap: &client_corev1.ConfigMapVolumeSourceApplyConfiguration{
								LocalObjectReferenceApplyConfiguration: client_corev1.LocalObjectReferenceApplyConfiguration{
									Name: pointer.String(option.Exe.Name),
								},
								Items: []client_corev1.KeyToPathApplyConfiguration{
									{
										Key:  pointer.String(option.ExeItemKey),
										Path: pointer.String("server"),
									},
								},
								DefaultMode: pointer.Int32(0777),
							},
						},
					},
					{
						// dataCollection服务运行所需的配置文件volume，包括通用配置
						Name: pointer.String("config"),
						VolumeSourceApplyConfiguration: client_corev1.VolumeSourceApplyConfiguration{
							// 通用配置
							ConfigMap: &client_corev1.ConfigMapVolumeSourceApplyConfiguration{
								LocalObjectReferenceApplyConfiguration: client_corev1.LocalObjectReferenceApplyConfiguration{
									Name: pointer.String("config"),
								},
								Items: []client_corev1.KeyToPathApplyConfiguration{
									{
										Key:  pointer.String("data-collection-config.yaml"),
										Path: pointer.String("config.yaml"),
									},
								},
							},
						},
					},
				},
				Containers: []client_corev1.ContainerApplyConfiguration{
					{
						Name:            &name,
						Image:           &option.Image,
						ImagePullPolicy: &imagePullPolicy,
						VolumeMounts: []client_corev1.VolumeMountApplyConfiguration{
							// 挂载放置编译后得到的二进制可执行程序的共享目录
							{
								Name:      pointer.String("app"),
								MountPath: pointer.String("/app"),
							},
							// 挂载应用运行所需的配置文件目录
							{
								Name:      pointer.String("config"),
								MountPath: pointer.String("/etc/app-configs"),
							},
						},
						Env: []client_corev1.EnvVarApplyConfiguration{
							{
								Name:  pointer.String("USERNAME"),
								Value: pointer.String(option.Username),
							},
							{
								// pod用于向kong网关注册服务时使用的名称
								Name: pointer.String("SERVICE_NAME"),
								ValueFrom: &client_corev1.EnvVarSourceApplyConfiguration{
									FieldRef: &client_corev1.ObjectFieldSelectorApplyConfiguration{
										FieldPath: pointer.String("metadata.name"),
									},
								},
							},
							{
								// statefulSet使用的headless service的名称，用于给pod组建域名
								Name:  pointer.String("HEADLESS_SERVICE_NAME"),
								Value: pointer.String(headlessServiceName),
							},
							{
								// 项目使用的域名，用于服务向网关注册路由时，增加Host请求头的匹配规则
								Name:  pointer.String("APP_DOMAIN_NAME"),
								Value: pointer.String(option.AppDomainName),
							},
						},
						Ports: []client_corev1.ContainerPortApplyConfiguration{
							{
								ContainerPort: pointer.Int32(8000),
							},
							{
								ContainerPort: pointer.Int32(9000),
							},
						},
						LivenessProbe: &client_corev1.ProbeApplyConfiguration{
							HandlerApplyConfiguration: client_corev1.HandlerApplyConfiguration{
								TCPSocket: &client_corev1.TCPSocketActionApplyConfiguration{
									Port: &servicePort,
								},
							},
							InitialDelaySeconds: pointer.Int32(20),
							PeriodSeconds:       pointer.Int32(60),
						},
						ReadinessProbe: &client_corev1.ProbeApplyConfiguration{
							HandlerApplyConfiguration: client_corev1.HandlerApplyConfiguration{
								TCPSocket: &client_corev1.TCPSocketActionApplyConfiguration{
									Port: &servicePort,
								},
							},
							InitialDelaySeconds: pointer.Int32(20),
							PeriodSeconds:       pointer.Int32(20),
						},
					},
				},
				RestartPolicy: &restartPolicy,
			},
		},
		// 使用的无头服务名
		ServiceName: &headlessServiceName,
	}
}
