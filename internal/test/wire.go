//go:build wireinject
// +build wireinject

// The build tag makes sure the stub is not built in the final build.

package test

import (
	"gitee.com/moyusir/service-centre/internal/biz"
	"gitee.com/moyusir/service-centre/internal/conf"
	"gitee.com/moyusir/service-centre/internal/data"
	"gitee.com/moyusir/service-centre/internal/server"
	"gitee.com/moyusir/service-centre/internal/service"
	"github.com/go-kratos/kratos/v2"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/google/wire"
)

// initApp init kratos application.
func initApp(*conf.Server, *conf.Service, *conf.Data, log.Logger) (*kratos.App, func(), error) {
	panic(wire.Build(server.ProviderSet, data.ProviderSet, biz.ProviderSet, service.ProviderSet, newApp))
}
