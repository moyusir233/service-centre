package biz

import (
	v1 "gitee.com/moyusir/service-centre/api/serviceCenter/v1"
	"gitee.com/moyusir/service-centre/internal/biz/codegenerator"
	"gitee.com/moyusir/service-centre/internal/biz/gateway"
	"gitee.com/moyusir/service-centre/internal/biz/kubecontroller"
	"gitee.com/moyusir/service-centre/internal/conf"
	"github.com/go-kratos/kratos/v2/errors"
	"github.com/go-kratos/kratos/v2/log"
	"golang.org/x/sync/errgroup"
	corev1 "k8s.io/api/core/v1"
	"time"
)

type UserUsecase struct {
	repo       UserRepo
	generator  *codegenerator.CodeGenerator
	controller *kubecontroller.KubeController
	gateway    *gateway.Manager
	logger     *log.Helper
}
type UserRepo interface {
	// Login 用户登录
	Login(username, password string) (token string, err error)
	// Register 用户注册
	Register(username, password, token string) error
	// UnRegister 用户注销
	UnRegister(username string) error
}

func NewUserUsecase(server *conf.Server, service *conf.Service, repo UserRepo, logger log.Logger) (*UserUsecase, error) {
	generator, err := codegenerator.NewCodeGenerator(
		service.CodeGenerator.DataProcessingTmplRoot,
		service.CodeGenerator.DataCollectionTmplRoot,
	)
	if err != nil {
		return nil, err
	}

	controller, err := kubecontroller.NewKubeController(server.Cluster.Namespace)
	if err != nil {
		return nil, err
	}

	manager, err := gateway.NewManager(server.Gateway.Address, server.AppDomainName)
	if err != nil {
		return nil, err
	}

	return &UserUsecase{
		repo:       repo,
		generator:  generator,
		controller: controller,
		gateway:    manager,
		logger:     log.NewHelper(logger),
	}, nil
}

func (u *UserUsecase) Login(username, password string) (token string, err error) {
	return u.repo.Login(username, password)
}

// Register 完成用户注册
// 首先向网关创建相应consumer，获得apiKey，然后在数据库中保存用户信息，然后进行生成服务代码，
// 然后为服务代码创建相应的configMap，然后启动相应的service和deployment，
// 最后在网关创建服务、路由以及认证插件
func (u *UserUsecase) Register(request *v1.RegisterRequest) (token string, err error) {
	if request == nil {
		return "", errors.BadRequest("request is nil", "")
	}
	username := request.User.Id

	defer func() {
		if err != nil {
			// 注册失败时，需要清理掉创建的无效资源
			u.Unregister(username)
		}
	}()
	// 在网关创建用户对应的consumer，从而创建token
	token, err = u.gateway.CreateConsumerAndKey(username)
	if err != nil {
		return "", err
	}

	// 往数据库中保存用户信息
	err = u.repo.Register(username, request.User.Password, token)
	if err != nil {
		return "", err
	}

	// 依据用户注册信息生成代码
	dcCode, dpCode, err := u.generator.GetServiceFiles(request.DeviceConfigRegisterInfos, request.DeviceStateRegisterInfos)
	if err != nil {
		return "", err
	}

	// 为用户创建服务运行所需的k8s资源

	// 创建代码对应的configMap
	dcCm, dpCm, err := u.controller.CreateConfigMapOfGeneratedCode(username, dcCode, dpCode)
	if err != nil {
		return "", err
	}

	// 创建注册信息对应的configMap
	registerInfo, err := u.controller.CreateConfigMapOfStateRegisterInfo(username, request.DeviceStateRegisterInfos)
	if err != nil {
		return "", err
	}

	// 利用协程部署数据收集和数据处理服务
	var dcService, dpService *corev1.Service
	eg := &errgroup.Group{}
	eg.Go(func() error {
		dcService, err = u.controller.DeployDataCollectionService(&kubecontroller.DataCollectionDeployOption{
			BaseDeployOption: kubecontroller.BaseDeployOption{
				Username:           username,
				Replica:            2,
				Timeout:            5 * time.Minute,
				ProjectRepoAddress: "git@gitee.com:moyusir/data-collection.git",
				ProjectBranch:      "code-template",
				ProjectDir:         "data-collection",
				ProjectApiDir:      "api/dataCollection/v1",
				ProjectServiceDir:  "internal/service",
				Image:              "moyusir233/graduation-design:data-collection",
			},
			Code:          dcCm,
			AppDomainName: u.gateway.AppDomainName,
		})
		return err
	})
	eg.Go(func() error {
		dpService, err = u.controller.DeployDataProcessingService(&kubecontroller.DataProcessingDeployOption{
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
		return err
	})
	err = eg.Wait()
	if err != nil {
		return "", err
	}

	// 为部署的服务向网关创建外部的路由
	err = u.gateway.CreateDcServiceRoute(username, dcService)
	if err != nil {
		return "", err
	}
	err = u.gateway.CreateDpServiceRoute(username, dpService)
	if err != nil {
		return "", err
	}

	return
}

// Unregister 注销用户，清理用户相关的资源，包括网关组件、k8s资源以及数据库的记录
func (u *UserUsecase) Unregister(username string) {
	// TODO 错误处理
	u.repo.UnRegister(username)
	u.gateway.Unregister(username)
	u.controller.Unregister(username)
}
