package biz

import (
	v1 "gitee.com/moyusir/service-centre/api/serviceCenter/v1"
	"gitee.com/moyusir/service-centre/internal/biz/gateway"
	"gitee.com/moyusir/service-centre/internal/biz/influxdb"
	"gitee.com/moyusir/service-centre/internal/biz/kubecontroller"
	"gitee.com/moyusir/service-centre/internal/conf"
	"github.com/go-kratos/kratos/v2/errors"
	"github.com/go-kratos/kratos/v2/log"
	"golang.org/x/sync/errgroup"
	corev1 "k8s.io/api/core/v1"
	"time"
)

type UserUsecase struct {
	repo                     UserRepo
	controller               *kubecontroller.KubeController
	gateway                  *gateway.Manager
	compilationCenterAddress string
	influxdbClient           *influxdb.Client
	logger                   *log.Helper
}
type UserRepo interface {
	// Login 用户登录
	Login(username, password string) (token string, err error)
	// Register 用户注册
	Register(username, password, token string) error
	// UnRegister 用户注销
	UnRegister(username string) error
	// GetClientCode 获得生成的客户端代码
	GetClientCode(username string) ([]byte, error)
}

func NewUserUsecase(server *conf.Server, repo UserRepo, logger log.Logger) (*UserUsecase, error) {
	controller, err := kubecontroller.NewKubeController(server.Cluster.Namespace)
	if err != nil {
		return nil, err
	}

	manager, err := gateway.NewManager(server.Gateway.Address, server.AppDomainName)
	if err != nil {
		return nil, err
	}

	influxdbClient, err := influxdb.NewInfluxdbClient(
		server.Influxdb.ServerUrl, server.Influxdb.AuthToken, server.Influxdb.Org)
	if err != nil {
		return nil, err
	}

	return &UserUsecase{
		repo:                     repo,
		controller:               controller,
		gateway:                  manager,
		compilationCenterAddress: server.CompilationCenter.Address,
		influxdbClient:           influxdbClient,
		logger:                   log.NewHelper(logger),
	}, nil
}

func (u *UserUsecase) Login(username, password string) (token string, err error) {
	return u.repo.Login(username, password)
}

// Register 完成用户注册
// 首先向网关创建相应consumer，获得apiKey，然后在数据库中保存用户信息，然后为用户创建相应的bucket，
// 然后进行生成服务代码，为服务代码创建相应的configMap，然后启动相应的service和deployment，
// 最后在网关创建服务、路由以及认证插件
func (u *UserUsecase) Register(request *v1.RegisterRequest) (token string, err error) {
	if request == nil {
		return "", errors.BadRequest("request is nil", "")
	}
	username := request.User.Id

	defer func() {
		if err != nil {
			// 注册失败时，需要清理掉创建的无效资源
			u.clear(username)
		}
	}()
	// 在网关创建用户对应的consumer，从而创建token
	token, err = u.gateway.CreateConsumerAndKey(username)
	if err != nil {
		return "", errors.Newf(
			500, "Register_Error",
			"创建用户对应的网关consumer时发生了错误:%v", err,
		)
	}

	// 往数据库中保存用户信息
	err = u.repo.Register(username, request.User.Password, token)
	if err != nil {
		return "", err
	}

	// 创建保存用户设备状态信息的influxdb bucket
	err = u.influxdbClient.CreateBucket(username)
	if err != nil {
		return "", errors.Newf(
			500, "Register_Error",
			"创建用户对应的influx bucket时发生了错误:%v", err,
		)
	}

	// 为用户创建服务运行所需的k8s资源
	// 创建用户注册信息对应的configMap，用于初始容器向编译中心发起编译请求使用
	registerInfo, err := u.controller.CreateConfigMapOfRegisterInfo(
		username, request.DeviceStateRegisterInfos, request.DeviceConfigRegisterInfos)
	if err != nil {
		return "", errors.Newf(
			500, "Register_Error",
			"创建用户对应的k8s资源时发生了错误:%v", err,
		)
	}

	// 利用协程部署数据收集和数据处理服务
	var dcService, dpService *corev1.Service
	eg := &errgroup.Group{}
	eg.Go(func() error {
		dcService, err = u.controller.DeployDataCollectionService(&kubecontroller.DataCollectionDeployOption{
			BaseDeployOption: kubecontroller.BaseDeployOption{
				Username:                 username,
				Replica:                  2,
				Timeout:                  5 * time.Minute,
				CompilationCenterAddress: u.compilationCenterAddress,
				RegisterInfo:             registerInfo,
				Image:                    "moyusir233/graduation-design:data-collection",
			},
			AppDomainName: u.gateway.AppDomainName,
		})
		return err
	})
	eg.Go(func() error {
		dpService, err = u.controller.DeployDataProcessingService(&kubecontroller.DataProcessingDeployOption{
			BaseDeployOption: kubecontroller.BaseDeployOption{
				Username:                 username,
				Replica:                  1,
				Timeout:                  5 * time.Minute,
				CompilationCenterAddress: u.compilationCenterAddress,
				RegisterInfo:             registerInfo,
				Image:                    "moyusir233/graduation-design:data-processing",
			},
		})
		return err
	})
	err = eg.Wait()
	if err != nil {
		return "", errors.Newf(
			500, "Register_Error",
			"创建用户服务相应的运行容器时发生了错误:%v", err,
		)
	}

	// 为部署的服务向网关创建外部的路由
	err = u.gateway.CreateDcServiceRoute(username, dcService)
	if err != nil {
		return "", errors.Newf(
			500, "Register_Error",
			"创建用户服务相应的路由时发生了错误:%v", err,
		)
	}
	err = u.gateway.CreateDpServiceRoute(username, dpService)
	if err != nil {
		return "", errors.Newf(
			500, "Register_Error",
			"创建用户服务相应的路由时发生了错误:%v", err,
		)
	}

	return
}

// Unregister 注销用户，清理用户相关的资源，包括网关组件、k8s资源以及数据库的记录
func (u *UserUsecase) Unregister(username, password string) error {
	// 确认密码是否正确，确保是用户本人操作的注销
	_, err := u.repo.Login(username, password)
	if err != nil {
		return errors.Forbidden(
			"UnRegister_Error", "账号或密码错误，无法执行注销操作")
	}

	return u.clear(username)
}

func (u *UserUsecase) GetClientCode(username string) ([]byte, error) {
	return u.repo.GetClientCode(username)
}

// 清理用户相关的资源
func (u *UserUsecase) clear(username string) (err error) {
	// TODO 错误处理
	err = u.repo.UnRegister(username)
	err = u.gateway.Unregister(username)
	err = u.influxdbClient.ClearBucket(username)
	err = u.controller.Unregister(username)

	return err
}
