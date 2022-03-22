package biz

import (
	"bytes"
	"context"
	compileApi "gitee.com/moyusir/compilation-center/api/compilationCenter/v1"
	v1 "gitee.com/moyusir/service-centre/api/serviceCenter/v1"
	"gitee.com/moyusir/service-centre/internal/biz/gateway"
	"gitee.com/moyusir/service-centre/internal/biz/kubecontroller"
	"gitee.com/moyusir/service-centre/internal/conf"
	"github.com/go-kratos/kratos/v2/errors"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/transport/grpc"
	"golang.org/x/sync/errgroup"
	corev1 "k8s.io/api/core/v1"
	"time"
)

type UserUsecase struct {
	repo       UserRepo
	controller *kubecontroller.KubeController
	gateway    *gateway.Manager
	client     compileApi.BuildClient
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

func NewUserUsecase(server *conf.Server, repo UserRepo, logger log.Logger) (*UserUsecase, func(), error) {
	controller, err := kubecontroller.NewKubeController(server.Cluster.Namespace)
	if err != nil {
		return nil, nil, err
	}

	manager, err := gateway.NewManager(server.Gateway.Address, server.AppDomainName)
	if err != nil {
		return nil, nil, err
	}

	conn, err := grpc.DialInsecure(context.Background(), grpc.WithEndpoint(server.CompilationCenter.Address))
	if err != nil {
		return nil, nil, err
	}
	client := compileApi.NewBuildClient(conn)

	return &UserUsecase{
		repo:       repo,
		controller: controller,
		gateway:    manager,
		client:     client,
		logger:     log.NewHelper(logger),
	}, func() { conn.Close() }, nil
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
			u.clear(username)
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

	// 调用编译中心的grpc编译服务，获得传输可执行程序二进制信息的grpc流，并读取二进制信息
	stream, err := u.client.GetServiceProgram(context.Background(), &compileApi.BuildRequest{
		Username:                  username,
		DeviceStateRegisterInfos:  request.DeviceStateRegisterInfos,
		DeviceConfigRegisterInfos: request.DeviceConfigRegisterInfos,
	})
	if err != nil {
		return "", err
	}

	// 读取流中的二进制数据
	dcExe := bytes.NewBuffer(make([]byte, 0, 1024))
	dpExe := bytes.NewBuffer(make([]byte, 0, 1024))

	for {
		reply, err := stream.Recv()
		if err != nil {
			break
		}
		if len(reply.DcExe) > 0 {
			dcExe.Write(reply.DcExe)
		}
		if len(reply.DpExe) > 0 {
			dpExe.Write(reply.DpExe)
		}
	}

	// 为用户创建服务运行所需的k8s资源

	// 创建可执行程序对应的configMap
	configMapOfExe, err := u.controller.CreateConfigMapOfExe(username, map[string][]byte{
		"dc": dcExe.Bytes(),
		"dp": dpExe.Bytes(),
	})
	if err != nil {
		return "", err
	}

	// 创建用户设备状态注册信息对应的configMap
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
				Username:   username,
				Replica:    2,
				Timeout:    5 * time.Minute,
				Image:      "moyusir233/graduation-design:data-collection",
				Exe:        configMapOfExe,
				ExeItemKey: "dc",
			},
			AppDomainName: u.gateway.AppDomainName,
		})
		return err
	})
	eg.Go(func() error {
		dpService, err = u.controller.DeployDataProcessingService(&kubecontroller.DataProcessingDeployOption{
			BaseDeployOption: kubecontroller.BaseDeployOption{
				Username:   username,
				Replica:    1,
				Timeout:    5 * time.Minute,
				Image:      "moyusir233/graduation-design:data-processing",
				Exe:        configMapOfExe,
				ExeItemKey: "dp",
			},
			RegisterInfo: registerInfo,
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
func (u *UserUsecase) Unregister(username, password string) error {
	// 确认密码是否正确，确保是用户本人操作的注销
	_, err := u.repo.Login(username, password)
	if err != nil {
		return errors.Forbidden("", "")
	}

	return u.clear(username)
}

// 清理用户相关的资源
func (u *UserUsecase) clear(username string) (err error) {
	// TODO 错误处理
	err = u.repo.UnRegister(username)
	err = u.gateway.Unregister(username)
	err = u.controller.Unregister(username)

	return err
}
