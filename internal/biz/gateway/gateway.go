package gateway

import (
	"gitee.com/moyusir/util/kong"
	"github.com/go-kratos/kratos/v2/errors"
	corev1 "k8s.io/api/core/v1"
	"net/http"
	"strings"
)

type Manager struct {
	*kong.Admin
	AppDomainName string
}

func NewManager(address, appDomainName string) (*Manager, error) {
	admin, err := kong.NewAdmin(address)
	if err != nil {
		return nil, err
	}
	return &Manager{Admin: admin, AppDomainName: appDomainName}, nil
}

// CreateConsumerAndKey 为用户创建在网关中的consumer实体以及相应的api密钥
func (m *Manager) CreateConsumerAndKey(username string) (apiKey string, err error) {
	consumerCreateOption := &kong.ConsumerCreateOption{
		Username: username,
		Tags:     []string{username},
	}
	consumer, err := m.Create(consumerCreateOption)
	if err != nil {
		return "", err
	}

	keyCreateOption := &kong.KeyCreateOption{Username: username}
	key, err := m.Create(keyCreateOption)
	if err != nil {
		m.Delete(consumer)
		return "", err
	}

	return key.(*kong.Key).Key, nil
}

// Unregister 清空用户在网关相关的组件
func (m *Manager) Unregister(username string) error {
	m.Clear(
		kong.FLAG_SERVICE|kong.FLAG_ROUTE|kong.FLAG_PLUGIN|kong.FLAG_CONSUMER, username)
	return nil
}

// CreateDcServiceRoute 为数据收集服务的service组件创建外部路由
func (m *Manager) CreateDcServiceRoute(username string, service *corev1.Service) error {
	if service == nil {
		return errors.New(500, "service is nil", "")
	}
	// 为数据收集服务的grpc连接创建路由
	// 查询service提供grpc的端口，默认为9000
	var grpcPort int32 = 9000
	var httpPort int32 = 8000
	for _, p := range service.Spec.Ports {
		if p.Name == "grpc" {
			grpcPort = p.Port
		} else if p.Name == "http" {
			httpPort = p.Port
		}
	}

	var (
		err     error
		objects []kong.Object
	)
	defer func() {
		if err != nil {
			for _, o := range objects {
				m.Admin.Delete(o)
			}
		}
	}()

	// 配置kong service组件的创建选项，需要附上用户名的tag方便后续用户注销
	serviceCreateOption := &kong.ServiceCreateOption{
		Name:     service.Name,
		Protocol: "grpc",
		// k8s中服务名即相应的域名
		Host:           service.Name,
		Port:           int(grpcPort),
		Enabled:        true,
		WriteTimeout:   600000,
		ReadTimeout:    600000,
		ConnectTimeout: 600000,
		Tags:           []string{username},
	}
	svc, err := m.Create(serviceCreateOption)
	if err != nil {
		return err
	}
	objects = append(objects, svc)

	configUpdateSvcName := service.Name + "-config-update"
	configUpdateServiceCreateOption := &kong.ServiceCreateOption{
		Name:     configUpdateSvcName,
		Protocol: "http",
		// k8s中服务名即相应的域名
		Host:    service.Name,
		Port:    int(httpPort),
		Path:    "/",
		Enabled: true,
		Tags:    []string{username},
	}
	configUpdateSvc, err := m.Create(configUpdateServiceCreateOption)
	if err != nil {
		return err
	}
	objects = append(objects, configUpdateSvc)

	// 创建路由，路由匹配条件包括host请求头和X-Service-Type:<用户名>-dc，
	// tag部分需要附上用户名，方便后续用户注销
	routeCreateOption := &kong.RouteCreateOption{
		Name:      service.Name,
		Protocols: []string{"grpc"},
		Hosts:     []string{m.AppDomainName},
		Paths:     []string{"/"},
		Headers: map[string][]string{
			"X-Service-Type": {username + "-dc"},
		},
		Service: &struct {
			Name string `json:"name,omitempty"`
			Id   string `json:"id,omitempty"`
		}{Name: service.Name},
		Tags: []string{username},
	}
	route, err := m.Create(routeCreateOption)
	if err != nil {
		return err
	}
	objects = append(objects, route)

	configUpdateRouteCreateOption := &kong.RouteCreateOption{
		Name:      configUpdateSvcName,
		Protocols: []string{"http"},
		Methods:   []string{http.MethodPost},
		Hosts:     []string{m.AppDomainName},
		Paths:     []string{"/"},
		Headers: map[string][]string{
			"X-Service-Type": {username + "-dc-config-update"},
		},
		StripPath: false,
		Service: &struct {
			Name string `json:"name,omitempty"`
			Id   string `json:"id,omitempty"`
		}{Name: configUpdateSvcName},
		Tags: []string{username},
	}
	configUpdateRoute, err := m.Create(configUpdateRouteCreateOption)
	if err != nil {
		return err
	}
	objects = append(objects, configUpdateRoute)

	// 创建认证插件，要求外部的grpc请求在query或者header中添加注册时得到的X-Api-Key
	pluginCreateOption := &kong.KeyAuthPluginCreateOption{
		Enabled: true,
		Service: &struct {
			Name string `json:"name,omitempty"`
			Id   string `json:"id,omitempty"`
		}{Name: service.Name},
		Config: &kong.KeyAuthPluginConfig{
			KeyNames:    []string{"X-Api-Key"},
			KeyInQuery:  true,
			KeyInHeader: true,
			KeyInBody:   false,
		},
	}
	plugin, err := m.Create(pluginCreateOption)
	if err != nil {
		return err
	}
	objects = append(objects, plugin)

	pluginCreateOption2 := &kong.KeyAuthPluginCreateOption{
		Enabled: true,
		Service: &struct {
			Name string `json:"name,omitempty"`
			Id   string `json:"id,omitempty"`
		}{Name: configUpdateSvcName},
		Config: &kong.KeyAuthPluginConfig{
			KeyNames:    []string{"X-Api-Key"},
			KeyInQuery:  true,
			KeyInHeader: true,
			KeyInBody:   false,
		},
	}
	_, err = m.Create(pluginCreateOption2)
	if err != nil {
		return err
	}

	return nil
}

// CreateDpServiceRoute 为数据处理服务的service组件创建外部路由
func (m *Manager) CreateDpServiceRoute(username string, service *corev1.Service) error {
	if service == nil {
		return errors.New(500, "service is nil", "")
	}

	// 为数据处理服务的http连接创建路由
	// 查询service提供http服务的端口，默认为8000
	var port int32 = 8000
	for _, p := range service.Spec.Ports {
		if p.Name == "http" {
			port = p.Port
			break
		}
	}

	var (
		objects []kong.Object
		err     error
	)
	defer func() {
		if err != nil {
			for _, o := range objects {
				m.Delete(o)
			}
		}
	}()

	// 配置kong service组件的创建选项，需要附上用户名的tag方便后续用户注销
	serviceCreateOption := &kong.ServiceCreateOption{
		Name:     service.Name,
		Protocol: "http",
		// k8s中服务名即相应的域名
		Host:         service.Name,
		Port:         int(port),
		Path:         "/",
		WriteTimeout: 600000,
		ReadTimeout:  600000,
		Enabled:      true,
		Tags:         []string{username},
	}
	svc, err := m.Create(serviceCreateOption)
	if err != nil {
		return err
	}
	objects = append(objects, svc)

	// 创建路由，路由匹配条件包括host请求头和X-Service-Type:<用户名>-dp，
	// tag部分需要附上用户名，方便后续用户注销
	routeCreateOption := &kong.RouteCreateOption{
		Name:      service.Name,
		Protocols: []string{"http"},
		Methods:   []string{http.MethodGet, http.MethodPut, http.MethodDelete},
		Hosts:     []string{m.AppDomainName},
		Paths:     []string{"/"},
		Headers: map[string][]string{
			"X-Service-Type": {username + "-dp"},
		},
		StripPath: false,
		Service: &struct {
			Name string `json:"name,omitempty"`
			Id   string `json:"id,omitempty"`
		}{Name: service.Name},
		Tags: []string{username},
	}
	route, err := m.Create(routeCreateOption)
	if err != nil {
		return err
	}
	objects = append(objects, route)

	// 由于浏览器发起ws连接时无法添加请求头，
	// 因此需要为建立预警推送ws连接的服务额外增加一个基于path匹配的路由
	wsRouteCreateOption := &kong.RouteCreateOption{
		Name:      service.Name + "-warning-push",
		Protocols: []string{"http"},
		Methods:   []string{http.MethodGet},
		Hosts:     []string{m.AppDomainName},
		Paths:     []string{"/warnings/push/" + strings.Replace(username, "_", "-", -1)},
		StripPath: false,
		Service: &struct {
			Name string `json:"name,omitempty"`
			Id   string `json:"id,omitempty"`
		}{Name: service.Name},
		Tags: []string{username},
	}
	wsRoute, err := m.Create(wsRouteCreateOption)
	if err != nil {
		return err
	}
	objects = append(objects, wsRoute)

	// 创建认证插件，要求外部的http请求在query或者header中添加注册时得到的X-Api-Key
	pluginCreateOption := &kong.KeyAuthPluginCreateOption{
		Enabled: true,
		Service: &struct {
			Name string `json:"name,omitempty"`
			Id   string `json:"id,omitempty"`
		}{Name: service.Name},
		Config: &kong.KeyAuthPluginConfig{
			KeyNames:    []string{"X-Api-Key"},
			KeyInQuery:  true,
			KeyInHeader: true,
			KeyInBody:   false,
		},
	}
	_, err = m.Create(pluginCreateOption)
	if err != nil {
		return err
	}

	return nil
}

// GetUsernameOfToken 获得与token相关的用户名
func (m *Manager) GetUsernameOfToken(token string) (string, error) {
	result := &struct {
		Username string `json:"username"`
	}{}

	response, err := m.Client.R().
		SetPathParam("token", token).
		SetResult(result).
		Get("/key-auths/{token}/consumer")
	if err != nil {
		return "", errors.Newf(500, "获得token相关的用户名时发生了错误: %s", err.Error())
	}
	if response.IsError() {
		return "", errors.Newf(
			500, "获得token相关的用户名时发生了错误: %s", response.String())
	}

	if result.Username == "" {
		return "", errors.Newf(
			400, "与该token相关的用户不存在: %s", token)
	}

	return result.Username, nil
}
