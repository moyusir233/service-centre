package kubecontroller

import (
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
	client_appsv1 "k8s.io/client-go/applyconfigurations/apps/v1"
	client_corev1 "k8s.io/client-go/applyconfigurations/core/v1"
	client_metav1 "k8s.io/client-go/applyconfigurations/meta/v1"
	"testing"
)

func Test_baseKubeController_CreateDeployment(t *testing.T) {
	controller, err := newBaseKubeController("test")
	if err != nil {
		t.Fatal(err)
	}

	var repicas int32 = 3
	name := "test"
	image := "nginx"
	pullpolicy := corev1.PullIfNotPresent
	port := intstr.FromInt(80)

	deploySpec := client_appsv1.DeploymentSpecApplyConfiguration{
		Replicas: &repicas,
		Selector: &client_metav1.LabelSelectorApplyConfiguration{
			MatchLabels: map[string]string{"app": "test"},
		},
		Template: &client_corev1.PodTemplateSpecApplyConfiguration{
			ObjectMetaApplyConfiguration: &client_metav1.ObjectMetaApplyConfiguration{
				Name:   &name,
				Labels: map[string]string{"app": "test"},
			},
			Spec: &client_corev1.PodSpecApplyConfiguration{
				Containers: []client_corev1.ContainerApplyConfiguration{
					{
						Name:            &name,
						Image:           &image,
						ImagePullPolicy: &pullpolicy,
						ReadinessProbe: &client_corev1.ProbeApplyConfiguration{
							HandlerApplyConfiguration: client_corev1.HandlerApplyConfiguration{
								HTTPGet: &client_corev1.HTTPGetActionApplyConfiguration{
									Port: &port,
								},
							},
						},
					},
				},
			},
		},
	}

	deployment, err := controller.CreateDeployment("test", map[string]string{"app": "test"}, &deploySpec)
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() {
		controller.DeleteResource("deployment", deployment.Name, deployment.TypeMeta)
	})
}

func Test_baseKubeController_CreateStatefulSet(t *testing.T) {
	controller, err := newBaseKubeController("test")
	if err != nil {
		t.Fatal(err)
	}

	// 为statefulSet创建一个无头服务
	var (
		portName          = "http"
		port        int32 = 80
		clusterIP         = "None"
		serviceType       = corev1.ServiceTypeClusterIP
		label             = map[string]string{"app": "test"}
	)
	serviceSpec := client_corev1.ServiceSpecApplyConfiguration{
		Ports: []client_corev1.ServicePortApplyConfiguration{
			{
				Name: &portName,
				Port: &port,
			},
		},
		Selector:  label,
		ClusterIP: &clusterIP,
		Type:      &serviceType,
	}
	service, err := controller.CreateService("test", label, &serviceSpec)
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() {
		controller.DeleteResource("service", service.Name, service.TypeMeta)
	})

	// 创建statefulSet
	var (
		replicas      int32 = 3
		name                = "test"
		image               = "nginx"
		pullpolicy          = corev1.PullIfNotPresent
		containerPort       = intstr.FromInt(80)
	)
	statefulSetSpec := client_appsv1.StatefulSetSpecApplyConfiguration{
		Replicas: &replicas,
		Selector: &client_metav1.LabelSelectorApplyConfiguration{
			MatchLabels: label,
		},
		Template: &client_corev1.PodTemplateSpecApplyConfiguration{
			ObjectMetaApplyConfiguration: &client_metav1.ObjectMetaApplyConfiguration{
				Labels: label,
			},
			Spec: &client_corev1.PodSpecApplyConfiguration{
				Containers: []client_corev1.ContainerApplyConfiguration{
					{
						Name:            &name,
						Image:           &image,
						ImagePullPolicy: &pullpolicy,
						ReadinessProbe: &client_corev1.ProbeApplyConfiguration{
							HandlerApplyConfiguration: client_corev1.HandlerApplyConfiguration{
								HTTPGet: &client_corev1.HTTPGetActionApplyConfiguration{
									Port: &containerPort,
								},
							},
						},
					},
				},
			},
		},
		ServiceName: &service.Name,
	}

	statefulSet, err := controller.CreateStatefulSet("test", label, &statefulSetSpec)
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() {
		controller.DeleteResource("statefulSet", statefulSet.Name, statefulSet.TypeMeta)
	})
}
