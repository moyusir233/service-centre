package test

import (
	"gitee.com/moyusir/service-centre/internal/biz/codegenerator"
	"gitee.com/moyusir/service-centre/internal/biz/kubecontroller"
	v1 "gitee.com/moyusir/util/api/util/v1"
	"golang.org/x/sync/errgroup"
	"google.golang.org/protobuf/types/known/durationpb"
	"testing"
	"time"
)

func TestKubeController_DeployService(t *testing.T) {
	var (
		username      = "test"
		appDomainName = "gd-k8s-master01"
	)

	generator, err := codegenerator.NewCodeGenerator(
		"/app/service-centre/internal/biz/codegenerator/data-processing-template",
		"/app/service-centre/internal/biz/codegenerator/data-collection-template")
	if err != nil {
		t.Fatal(err)
	}

	controller, err := kubecontroller.NewKubeController("test")
	if err != nil {
		t.Fatal(err)
	}

	// 测试将生成的代码进行编译和部署
	// 定义生成代码需要使用的注册信息
	configInfo := []*v1.DeviceConfigRegisterInfo{
		{
			Fields: []*v1.DeviceConfigRegisterInfo_Field{
				{
					Name: "id",
					Type: v1.Type_STRING,
				},
				{
					Name: "status",
					Type: v1.Type_BOOL,
				},
			},
		},
		{
			Fields: []*v1.DeviceConfigRegisterInfo_Field{
				{
					Name: "id",
					Type: v1.Type_STRING,
				},
				{
					Name: "status",
					Type: v1.Type_BOOL,
				},
			},
		},
	}
	stateInfo := []*v1.DeviceStateRegisterInfo{
		{
			Fields: []*v1.DeviceStateRegisterInfo_Field{
				{
					Name:        "id",
					Type:        v1.Type_STRING,
					WarningRule: nil,
				},
				{
					Name:        "time",
					Type:        v1.Type_TIMESTAMP,
					WarningRule: nil,
				},
				{
					Name: "current",
					Type: v1.Type_DOUBLE,
					WarningRule: &v1.DeviceStateRegisterInfo_WarningRule{
						CmpRule: &v1.DeviceStateRegisterInfo_CmpRule{
							Cmp: v1.DeviceStateRegisterInfo_GT,
							Arg: "1000",
						},
						AggregationOperation: v1.DeviceStateRegisterInfo_MIN,
						Duration:             durationpb.New(time.Minute),
					},
				},
				{
					Name: "voltage",
					Type: v1.Type_DOUBLE,
					WarningRule: &v1.DeviceStateRegisterInfo_WarningRule{
						CmpRule: &v1.DeviceStateRegisterInfo_CmpRule{
							Cmp: v1.DeviceStateRegisterInfo_GT,
							Arg: "1000",
						},
						AggregationOperation: v1.DeviceStateRegisterInfo_MIN,
						Duration:             durationpb.New(time.Minute),
					},
				},
			},
		},
		{
			Fields: []*v1.DeviceStateRegisterInfo_Field{
				{
					Name:        "id",
					Type:        v1.Type_STRING,
					WarningRule: nil,
				},
				{
					Name:        "time",
					Type:        v1.Type_TIMESTAMP,
					WarningRule: nil,
				},
				{
					Name: "current",
					Type: v1.Type_DOUBLE,
					WarningRule: &v1.DeviceStateRegisterInfo_WarningRule{
						CmpRule: &v1.DeviceStateRegisterInfo_CmpRule{
							Cmp: v1.DeviceStateRegisterInfo_GT,
							Arg: "1000",
						},
						AggregationOperation: v1.DeviceStateRegisterInfo_MIN,
						Duration:             durationpb.New(time.Minute),
					},
				},
				{
					Name: "voltage",
					Type: v1.Type_DOUBLE,
					WarningRule: &v1.DeviceStateRegisterInfo_WarningRule{
						CmpRule: &v1.DeviceStateRegisterInfo_CmpRule{
							Cmp: v1.DeviceStateRegisterInfo_GT,
							Arg: "1000",
						},
						AggregationOperation: v1.DeviceStateRegisterInfo_MIN,
						Duration:             durationpb.New(time.Minute),
					},
				},
			},
		},
	}

	// 获得生成的代码
	dcCode, dpCode, err := generator.GetServiceFiles(configInfo, stateInfo)
	if err != nil {
		t.Fatal(err)
	}

	// 为生成的代码创建configMap
	dcCm, dpCm, err := controller.CreateConfigMapOfGeneratedCode(username, dcCode, dpCode)
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() {
		t.Log(
			controller.DeleteResource(dpCm.Name, "ConfigMap"),
			controller.DeleteResource(dcCm.Name, "ConfigMap"),
		)
	})

	// 为用户注册信息创建configMap
	registerInfo, err := controller.CreateConfigMapOfStateRegisterInfo(username, stateInfo)
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() {
		t.Log(controller.DeleteResource(registerInfo.Name, "ConfigMap"))
	})

	// 然后部署服务
	eg := &errgroup.Group{}
	eg.Go(func() error {
		_, err := controller.DeployDataCollectionService(&kubecontroller.DataCollectionDeployOption{
			BaseDeployOption: kubecontroller.BaseDeployOption{
				Username:           username,
				Replica:            1,
				Timeout:            5 * time.Minute,
				ProjectRepoAddress: "git@gitee.com:moyusir/data-collection.git",
				ProjectBranch:      "code-template",
				ProjectDir:         "data-collection",
				ProjectApiDir:      "api/dataCollection/v1",
				ProjectServiceDir:  "internal/service",
				Image:              "moyusir233/graduation-design:data-collection",
			},
			Code:          dcCm,
			AppDomainName: appDomainName,
		})
		if err != nil {
			return err
		} else {
			return nil
		}
	})
	eg.Go(func() error {
		_, err := controller.DeployDataProcessingService(&kubecontroller.DataProcessingDeployOption{
			BaseDeployOption: kubecontroller.BaseDeployOption{
				Username:           username,
				Replica:            1,
				Timeout:            5 * time.Minute,
				ProjectRepoAddress: "git@gitee.com:moyusir/data-processing.git",
				ProjectBranch:      "code-template",
				ProjectDir:         "data-processing",
				ProjectApiDir:      "api/dataProcessing/v1",
				ProjectServiceDir:  "internal/service",
				Image:              "moyusir233/graduation-design:data-processing",
			},
			RegisterInfo: registerInfo,
			Code:         dpCm,
		})
		if err != nil {
			return err
		} else {
			return nil
		}
	})

	t.Cleanup(func() {
		var errs []error
		errs = append(errs,
			controller.DeleteResources("Deployment", "user="+username),
			controller.DeleteResources("StatefulSet", "user="+username),
			controller.DeleteResources("Service", "user="+username),
		)
		for _, e := range errs {
			t.Log(e)
		}
	})
	err = eg.Wait()
	if err != nil {
		t.Fatal(err)
	}
}
