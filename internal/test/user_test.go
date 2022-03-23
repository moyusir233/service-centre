package test

import (
	"context"
	"encoding/json"
	v1 "gitee.com/moyusir/service-centre/api/serviceCenter/v1"
	utilApi "gitee.com/moyusir/util/api/util/v1"
	"github.com/golang/protobuf/proto"
	"github.com/gorilla/websocket"
	"github.com/imroc/req/v3"
	"google.golang.org/protobuf/types/known/durationpb"
	"net/http"
	"testing"
	"time"
)

func TestUserUsecase(t *testing.T) {
	const (
		KONG_HTTP_URL = "http://kong.test.svc.cluster.local:8000"
	)
	// 定义生成代码需要使用的注册信息
	configInfo := []*utilApi.DeviceConfigRegisterInfo{
		{
			Fields: []*utilApi.DeviceConfigRegisterInfo_Field{
				{
					Name: "id",
					Type: utilApi.Type_STRING,
				},
				{
					Name: "status",
					Type: utilApi.Type_BOOL,
				},
			},
		},
		{
			Fields: []*utilApi.DeviceConfigRegisterInfo_Field{
				{
					Name: "id",
					Type: utilApi.Type_STRING,
				},
				{
					Name: "status",
					Type: utilApi.Type_BOOL,
				},
			},
		},
	}
	stateInfo := []*utilApi.DeviceStateRegisterInfo{
		{
			Fields: []*utilApi.DeviceStateRegisterInfo_Field{
				{
					Name:        "id",
					Type:        utilApi.Type_STRING,
					WarningRule: nil,
				},
				{
					Name:        "time",
					Type:        utilApi.Type_TIMESTAMP,
					WarningRule: nil,
				},
				{
					Name: "current",
					Type: utilApi.Type_INT64,
					WarningRule: &utilApi.DeviceStateRegisterInfo_WarningRule{
						CmpRule: &utilApi.DeviceStateRegisterInfo_CmpRule{
							Cmp: utilApi.DeviceStateRegisterInfo_GT,
							Arg: "1000",
						},
						AggregationOperation: utilApi.DeviceStateRegisterInfo_MIN,
						Duration:             durationpb.New(time.Minute),
					},
				},
				{
					Name: "voltage",
					Type: utilApi.Type_INT64,
					WarningRule: &utilApi.DeviceStateRegisterInfo_WarningRule{
						CmpRule: &utilApi.DeviceStateRegisterInfo_CmpRule{
							Cmp: utilApi.DeviceStateRegisterInfo_GT,
							Arg: "1000",
						},
						AggregationOperation: utilApi.DeviceStateRegisterInfo_MIN,
						Duration:             durationpb.New(time.Minute),
					},
				},
			},
		},
		{
			Fields: []*utilApi.DeviceStateRegisterInfo_Field{
				{
					Name:        "id",
					Type:        utilApi.Type_STRING,
					WarningRule: nil,
				},
				{
					Name:        "time",
					Type:        utilApi.Type_TIMESTAMP,
					WarningRule: nil,
				},
				{
					Name: "current",
					Type: utilApi.Type_DOUBLE,
					WarningRule: &utilApi.DeviceStateRegisterInfo_WarningRule{
						CmpRule: &utilApi.DeviceStateRegisterInfo_CmpRule{
							Cmp: utilApi.DeviceStateRegisterInfo_GT,
							Arg: "1000",
						},
						AggregationOperation: utilApi.DeviceStateRegisterInfo_MIN,
						Duration:             durationpb.New(time.Minute),
					},
				},
				{
					Name: "voltage",
					Type: utilApi.Type_DOUBLE,
					WarningRule: &utilApi.DeviceStateRegisterInfo_WarningRule{
						CmpRule: &utilApi.DeviceStateRegisterInfo_CmpRule{
							Cmp: utilApi.DeviceStateRegisterInfo_GT,
							Arg: "1000",
						},
						AggregationOperation: utilApi.DeviceStateRegisterInfo_MIN,
						Duration:             durationpb.New(time.Minute),
					},
				},
			},
		},
	}
	registerReq := &v1.RegisterRequest{
		User: &utilApi.User{
			Id:       "test",
			Password: "test",
		},
		DeviceConfigRegisterInfos: configInfo,
		DeviceStateRegisterInfos:  stateInfo,
	}

	// 注册，登录获取token，然后测试能否查询到设备状态注册信息、建立ws连接
	userHTTPClient := StartServiceCenterServer(t)
	_, err := userHTTPClient.Register(context.Background(), registerReq)
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() {
		_, err := userHTTPClient.Unregister(context.Background(), &utilApi.User{
			Id:       "test",
			Password: "test",
		})
		if err != nil {
			t.Error(err)
		}
	})

	reply, err := userHTTPClient.Login(context.Background(), &utilApi.User{
		Id:       "test",
		Password: "test",
	})
	if err != nil {
		t.Fatal(err)
	}

	// 发送查询设备状态注册信息的请求
	t.Run("Test_GetDeviceStateRegisterInfo", func(t *testing.T) {
		client := req.C().SetBaseURL(KONG_HTTP_URL)
		response, err := client.R().
			SetHeaders(map[string]string{
				"X-Api-Key":      reply.Token,
				"X-Service-Type": "test-dp",
			}).Get("/register-info/states/0")
		if err != nil {
			t.Fatal(err)
		}
		if response.IsError() {
			t.Fatal(response.Error())
		}

		registerInfo := new(utilApi.DeviceStateRegisterInfo)
		err = json.NewDecoder(response.Body).Decode(registerInfo)
		if err != nil {
			t.Fatal(err)
		}
		// 比较查询得到的注册信息和一开始注册上传的是否一样
		if !proto.Equal(registerReq, stateInfo[0]) {
			t.Error("wrong device state register info")
		}
	})

	t.Run("Test_WarningPushWebsocket", func(t *testing.T) {
		// 建立接收故障信息推送的ws连接
		conn, _, err := websocket.DefaultDialer.Dial(
			"ws://kong.test.svc.cluster.local:8000/warnings/push",
			http.Header{
				"X-Api-Key":      {reply.Token},
				"X-Service-Type": {"test-dp"},
			},
		)
		if err != nil {
			t.Error(err)
		}
		err = conn.Close()
		if err != nil {
			t.Error(err)
		}
	})
}
