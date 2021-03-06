// Code generated by protoc-gen-go-http. DO NOT EDIT.
// versions:
// protoc-gen-go-http v2.1.3

package v1

import (
	context "context"
	v1 "gitee.com/moyusir/util/api/util/v1"
	http "github.com/go-kratos/kratos/v2/transport/http"
	binding "github.com/go-kratos/kratos/v2/transport/http/binding"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the kratos package it is being compiled against.
var _ = new(context.Context)
var _ = binding.EncodeURL

const _ = http.SupportPackageIsVersion1

type UserHTTPServer interface {
	DownloadClientCode(context.Context, *DownloadClientCodeRequest) (*File, error)
	GetRegisterInfo(context.Context, *GetRegisterInfoRequest) (*GetRegisterInfoReply, error)
	Login(context.Context, *v1.User) (*LoginReply, error)
	Register(context.Context, *RegisterRequest) (*RegisterReply, error)
	Unregister(context.Context, *v1.User) (*UnregisterReply, error)
}

func RegisterUserHTTPServer(s *http.Server, srv UserHTTPServer) {
	r := s.Route("/")
	r.POST("/users", _User_Register0_HTTP_Handler(srv))
	r.GET("/users/register-info/{token}", _User_GetRegisterInfo0_HTTP_Handler(srv))
	r.GET("/users", _User_Login0_HTTP_Handler(srv))
	r.DELETE("/users", _User_Unregister0_HTTP_Handler(srv))
	r.GET("/users/client-code/{username}", _User_DownloadClientCode0_HTTP_Handler(srv))
}

func _User_Register0_HTTP_Handler(srv UserHTTPServer) func(ctx http.Context) error {
	return func(ctx http.Context) error {
		var in RegisterRequest
		if err := ctx.Bind(&in); err != nil {
			return err
		}
		http.SetOperation(ctx, "/api.serviceCentre.v1.User/Register")
		h := ctx.Middleware(func(ctx context.Context, req interface{}) (interface{}, error) {
			return srv.Register(ctx, req.(*RegisterRequest))
		})
		out, err := h(ctx, &in)
		if err != nil {
			return err
		}
		reply := out.(*RegisterReply)
		return ctx.Result(200, reply)
	}
}

func _User_GetRegisterInfo0_HTTP_Handler(srv UserHTTPServer) func(ctx http.Context) error {
	return func(ctx http.Context) error {
		var in GetRegisterInfoRequest
		if err := ctx.BindQuery(&in); err != nil {
			return err
		}
		if err := ctx.BindVars(&in); err != nil {
			return err
		}
		http.SetOperation(ctx, "/api.serviceCentre.v1.User/GetRegisterInfo")
		h := ctx.Middleware(func(ctx context.Context, req interface{}) (interface{}, error) {
			return srv.GetRegisterInfo(ctx, req.(*GetRegisterInfoRequest))
		})
		out, err := h(ctx, &in)
		if err != nil {
			return err
		}
		reply := out.(*GetRegisterInfoReply)
		return ctx.Result(200, reply)
	}
}

func _User_Login0_HTTP_Handler(srv UserHTTPServer) func(ctx http.Context) error {
	return func(ctx http.Context) error {
		var in v1.User
		if err := ctx.BindQuery(&in); err != nil {
			return err
		}
		http.SetOperation(ctx, "/api.serviceCentre.v1.User/Login")
		h := ctx.Middleware(func(ctx context.Context, req interface{}) (interface{}, error) {
			return srv.Login(ctx, req.(*v1.User))
		})
		out, err := h(ctx, &in)
		if err != nil {
			return err
		}
		reply := out.(*LoginReply)
		return ctx.Result(200, reply)
	}
}

func _User_Unregister0_HTTP_Handler(srv UserHTTPServer) func(ctx http.Context) error {
	return func(ctx http.Context) error {
		var in v1.User
		if err := ctx.BindQuery(&in); err != nil {
			return err
		}
		http.SetOperation(ctx, "/api.serviceCentre.v1.User/Unregister")
		h := ctx.Middleware(func(ctx context.Context, req interface{}) (interface{}, error) {
			return srv.Unregister(ctx, req.(*v1.User))
		})
		out, err := h(ctx, &in)
		if err != nil {
			return err
		}
		reply := out.(*UnregisterReply)
		return ctx.Result(200, reply)
	}
}

func _User_DownloadClientCode0_HTTP_Handler(srv UserHTTPServer) func(ctx http.Context) error {
	return func(ctx http.Context) error {
		var in DownloadClientCodeRequest
		if err := ctx.BindQuery(&in); err != nil {
			return err
		}
		if err := ctx.BindVars(&in); err != nil {
			return err
		}
		http.SetOperation(ctx, "/api.serviceCentre.v1.User/DownloadClientCode")
		h := ctx.Middleware(func(ctx context.Context, req interface{}) (interface{}, error) {
			return srv.DownloadClientCode(ctx, req.(*DownloadClientCodeRequest))
		})
		out, err := h(ctx, &in)
		if err != nil {
			return err
		}
		reply := out.(*File)
		return ctx.Result(200, reply)
	}
}

type UserHTTPClient interface {
	DownloadClientCode(ctx context.Context, req *DownloadClientCodeRequest, opts ...http.CallOption) (rsp *File, err error)
	GetRegisterInfo(ctx context.Context, req *GetRegisterInfoRequest, opts ...http.CallOption) (rsp *GetRegisterInfoReply, err error)
	Login(ctx context.Context, req *v1.User, opts ...http.CallOption) (rsp *LoginReply, err error)
	Register(ctx context.Context, req *RegisterRequest, opts ...http.CallOption) (rsp *RegisterReply, err error)
	Unregister(ctx context.Context, req *v1.User, opts ...http.CallOption) (rsp *UnregisterReply, err error)
}

type UserHTTPClientImpl struct {
	cc *http.Client
}

func NewUserHTTPClient(client *http.Client) UserHTTPClient {
	return &UserHTTPClientImpl{client}
}

func (c *UserHTTPClientImpl) DownloadClientCode(ctx context.Context, in *DownloadClientCodeRequest, opts ...http.CallOption) (*File, error) {
	var out File
	pattern := "/users/client-code/{username}"
	path := binding.EncodeURL(pattern, in, true)
	opts = append(opts, http.Operation("/api.serviceCentre.v1.User/DownloadClientCode"))
	opts = append(opts, http.PathTemplate(pattern))
	err := c.cc.Invoke(ctx, "GET", path, nil, &out, opts...)
	if err != nil {
		return nil, err
	}
	return &out, err
}

func (c *UserHTTPClientImpl) GetRegisterInfo(ctx context.Context, in *GetRegisterInfoRequest, opts ...http.CallOption) (*GetRegisterInfoReply, error) {
	var out GetRegisterInfoReply
	pattern := "/users/register-info/{token}"
	path := binding.EncodeURL(pattern, in, true)
	opts = append(opts, http.Operation("/api.serviceCentre.v1.User/GetRegisterInfo"))
	opts = append(opts, http.PathTemplate(pattern))
	err := c.cc.Invoke(ctx, "GET", path, nil, &out, opts...)
	if err != nil {
		return nil, err
	}
	return &out, err
}

func (c *UserHTTPClientImpl) Login(ctx context.Context, in *v1.User, opts ...http.CallOption) (*LoginReply, error) {
	var out LoginReply
	pattern := "/users"
	path := binding.EncodeURL(pattern, in, true)
	opts = append(opts, http.Operation("/api.serviceCentre.v1.User/Login"))
	opts = append(opts, http.PathTemplate(pattern))
	err := c.cc.Invoke(ctx, "GET", path, nil, &out, opts...)
	if err != nil {
		return nil, err
	}
	return &out, err
}

func (c *UserHTTPClientImpl) Register(ctx context.Context, in *RegisterRequest, opts ...http.CallOption) (*RegisterReply, error) {
	var out RegisterReply
	pattern := "/users"
	path := binding.EncodeURL(pattern, in, false)
	opts = append(opts, http.Operation("/api.serviceCentre.v1.User/Register"))
	opts = append(opts, http.PathTemplate(pattern))
	err := c.cc.Invoke(ctx, "POST", path, in, &out, opts...)
	if err != nil {
		return nil, err
	}
	return &out, err
}

func (c *UserHTTPClientImpl) Unregister(ctx context.Context, in *v1.User, opts ...http.CallOption) (*UnregisterReply, error) {
	var out UnregisterReply
	pattern := "/users"
	path := binding.EncodeURL(pattern, in, true)
	opts = append(opts, http.Operation("/api.serviceCentre.v1.User/Unregister"))
	opts = append(opts, http.PathTemplate(pattern))
	err := c.cc.Invoke(ctx, "DELETE", path, nil, &out, opts...)
	if err != nil {
		return nil, err
	}
	return &out, err
}
