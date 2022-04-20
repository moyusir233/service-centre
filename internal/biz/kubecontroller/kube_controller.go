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

	// 服务向编译中心发出请求时需要的配置

	// 编译中心的地址
	CompilationCenterAddress string
	// 以json形式保存用户注册信息的cm，其中状态信息以state.json为键，配置信息以config.json为键
	RegisterInfo *corev1.ConfigMap

	// 服务运行时配置
	// 运行服务使用的镜像
	Image string
}
type DataProcessingDeployOption struct {
	BaseDeployOption
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

func (c *KubeController) CreateConfigMapOfRegisterInfo(
	username string,
	states []*v1.DeviceStateRegisterInfo, configs []*v1.DeviceConfigRegisterInfo) (*corev1.ConfigMap, error) {
	// 以user:username为label,username-state-register-info为名称创建cm
	// 保存注册信息的json数据
	stateJson, err := json.Marshal(states)
	if err != nil {
		return nil, err
	}
	configJson, err := json.Marshal(configs)
	if err != nil {
		return nil, err
	}

	label := map[string]string{"user": username}
	return c.CreateConfigMap(username+"-register-info", label, map[string]string{
		"state.json":  string(stateJson),
		"config.json": string(configJson),
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
	// 解析json日志所需注解的前缀
	prefix := fmt.Sprintf("co.elastic.logs.%s/", name)

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
				Annotations: map[string]string{
					prefix + "json.keys_under_root": "true",
					prefix + "json.add_error_key":   "true",
					prefix + "json.message_key":     "msg",
					prefix + "json.expand_keys":     "true",
				},
			},
			Spec: &client_corev1.PodSpecApplyConfiguration{
				// 配置需要挂载的volume
				Volumes: []client_corev1.VolumeApplyConfiguration{
					{
						// 存放编译得到的二进制程序的中转目录
						Name: pointer.String("app"),
						VolumeSourceApplyConfiguration: client_corev1.VolumeSourceApplyConfiguration{
							EmptyDir: &client_corev1.EmptyDirVolumeSourceApplyConfiguration{},
						},
					},
					{
						// 存放用户注册信息的cm，用于向编译中心发起编译请求
						Name: pointer.String("register-info"),
						VolumeSourceApplyConfiguration: client_corev1.VolumeSourceApplyConfiguration{
							ConfigMap: &client_corev1.ConfigMapVolumeSourceApplyConfiguration{
								LocalObjectReferenceApplyConfiguration: client_corev1.LocalObjectReferenceApplyConfiguration{
									Name: &option.RegisterInfo.Name,
								},
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
											Items: []client_corev1.KeyToPathApplyConfiguration{
												{
													Key:  pointer.String("state.json"),
													Path: pointer.String("register_info.json"),
												},
											},
										},
									},
								},
							},
						},
					},
				},
				InitContainers: []client_corev1.ContainerApplyConfiguration{
					// 负责向编译中心发出编译请求，获得二进制可执行程序的initContainer
					{
						Name:  pointer.String("build"),
						Image: pointer.String("moyusir233/graduation-design:compilation-client"),
						Args: []string{
							"-u", option.Username,
							"-address", option.CompilationCenterAddress,
							"-state", "/etc/register-info/state.json",
							"-config", "/etc/register-info/config.json",
							"-service_type", "dp",
							"-o", "/app/server",
						},
						VolumeMounts: []client_corev1.VolumeMountApplyConfiguration{
							{
								// 用户注册信息
								Name:      pointer.String("register-info"),
								MountPath: pointer.String("/etc/register-info"),
							},
							{
								// 存放二进制文件的中转目录
								Name:      pointer.String("app"),
								MountPath: pointer.String("/app"),
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
	// 解析json日志所需注解的前缀
	prefix := fmt.Sprintf("co.elastic.logs.%s/", name)

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
				Annotations: map[string]string{
					prefix + "json.keys_under_root": "true",
					prefix + "json.add_error_key":   "true",
					prefix + "json.message_key":     "msg",
					prefix + "json.expand_keys":     "true",
				},
			},
			Spec: &client_corev1.PodSpecApplyConfiguration{
				// 配置需要挂载的volume
				Volumes: []client_corev1.VolumeApplyConfiguration{
					{
						// 存放编译得到的二进制程序的中转目录
						Name: pointer.String("app"),
						VolumeSourceApplyConfiguration: client_corev1.VolumeSourceApplyConfiguration{
							EmptyDir: &client_corev1.EmptyDirVolumeSourceApplyConfiguration{},
						},
					},
					{
						// 存放用户注册信息的cm，用于向编译中心发起编译请求
						Name: pointer.String("register-info"),
						VolumeSourceApplyConfiguration: client_corev1.VolumeSourceApplyConfiguration{
							ConfigMap: &client_corev1.ConfigMapVolumeSourceApplyConfiguration{
								LocalObjectReferenceApplyConfiguration: client_corev1.LocalObjectReferenceApplyConfiguration{
									Name: &option.RegisterInfo.Name,
								},
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
				InitContainers: []client_corev1.ContainerApplyConfiguration{
					// 负责向编译中心发出编译请求，获得二进制可执行程序的initContainer
					{
						Name:  pointer.String("build"),
						Image: pointer.String("moyusir233/graduation-design:compilation-client"),
						Args: []string{
							"-u", option.Username,
							"-address", option.CompilationCenterAddress,
							"-state", "/etc/register-info/state.json",
							"-config", "/etc/register-info/config.json",
							"-service_type", "dc",
							"-o", "/app/server",
						},
						VolumeMounts: []client_corev1.VolumeMountApplyConfiguration{
							{
								// 用户注册信息
								Name:      pointer.String("register-info"),
								MountPath: pointer.String("/etc/register-info"),
							},
							{
								// 存放二进制文件的中转目录
								Name:      pointer.String("app"),
								MountPath: pointer.String("/app"),
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
