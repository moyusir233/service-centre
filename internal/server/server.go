package server

import (
	"context"
	v1 "gitee.com/moyusir/service-centre/api/serviceCenter/v1"
	utilApi "gitee.com/moyusir/util/api/util/v1"
	"github.com/go-kratos/kratos/v2/errors"
	"github.com/go-kratos/kratos/v2/middleware"
	h "github.com/go-kratos/kratos/v2/transport/http"
	"github.com/google/wire"
	"net/http"
)

// ProviderSet is server providers.
var ProviderSet = wire.NewSet(NewHTTPServer)

// MyResponseEncoder 用于处理文件响应的响应编码器
func MyResponseEncoder(respWriter http.ResponseWriter, req *http.Request, v interface{}) error {
	if file, ok := v.(*v1.File); ok {
		// 设置文件下载的响应头
		respWriter.Header().Set("Content-Type", "application/octet-stream")
		respWriter.Header().Set("Content-Disposition", "attachment;filename="+file.Name)
		_, err := respWriter.Write(file.Content)
		return err
	} else {
		return h.DefaultResponseEncoder(respWriter, req, v)
	}
}

// RegisterValidator 用于验证用户注册信息的中间件.
func RegisterValidator() middleware.Middleware {
	return func(handler middleware.Handler) middleware.Handler {
		return func(ctx context.Context, req interface{}) (reply interface{}, err error) {
			if registerReq, ok := req.(*v1.RegisterRequest); ok {
				// 需要对注册信息进行三个方面的检测：
				// 1. 每台设备注册的字段名称不允许重复
				// 2. 不允许为非数值类型的字段注册预警规则
				// 3. 每台设备必须有类型为string，名为id的字段
				// 4. 在设备状态的注册字段中必须包含名为time，类别的timestamp的时间戳字段

				// 检查配置注册信息
				for _, config := range registerReq.DeviceConfigRegisterInfos {
					hasID := false
					fieldNames := make(map[string]bool, len(config.Fields))

					for _, field := range config.Fields {
						// 检查是否为id字段
						if !hasID && field.Name == "id" && field.Type == utilApi.Type_STRING {
							hasID = true
						}

						// 检查字段名是否重复
						if fieldNames[field.Name] {
							return nil, errors.BadRequest(
								"repeated field name",
								"The field name of a device cannot be duplicate")
						}
						fieldNames[field.Name] = true
					}

					if !hasID {
						return nil, errors.BadRequest(
							"missing id field",
							"The id field is missing")
					}
				}

				// 检查设备状态注册信息
				nonNumericType := map[utilApi.Type]bool{
					utilApi.Type_STRING:    true,
					utilApi.Type_TIMESTAMP: true,
					utilApi.Type_BOOL:      true,
				}
				for _, state := range registerReq.DeviceStateRegisterInfos {
					hasID := false
					hasTime := false
					fieldNames := make(map[string]bool, len(state.Fields))

					for _, field := range state.Fields {
						// 检查是否为id字段
						if !hasID && field.Name == "id" && field.Type == utilApi.Type_STRING {
							hasID = true
						}
						// 检查是否为time字段
						if !hasTime && field.Name == "time" && field.Type == utilApi.Type_TIMESTAMP {
							hasTime = true
						}

						// 检查字段名是否重复
						if fieldNames[field.Name] {
							return nil, errors.BadRequest(
								"repeated field name",
								"The field name of a device cannot be duplicate")
						}
						fieldNames[field.Name] = true

						// 检查是否有为非数值类型注册的预警规则
						if field.WarningRule != nil && nonNumericType[field.Type] {
							return nil, errors.BadRequest(
								"wrong warning rule",
								"warning rule cannot be registered for non numeric types")
						}
					}
					if !hasID {
						return nil, errors.BadRequest(
							"missing id field",
							"The id field is missing")
					}
					if !hasTime {
						return nil, errors.BadRequest(
							"missing time field",
							"The time field is missing")
					}
				}
			}
			return handler(ctx, req)
		}
	}
}
