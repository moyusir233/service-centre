package server

import (
	v1 "gitee.com/moyusir/service-centre/api/serviceCenter/v1"
	h "github.com/go-kratos/kratos/v2/transport/http"
	"github.com/google/wire"
	"net/http"
)

// ProviderSet is server providers.
var ProviderSet = wire.NewSet(NewHTTPServer)

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
