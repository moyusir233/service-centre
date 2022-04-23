package test

import (
	"context"
	v1 "gitee.com/moyusir/service-centre/api/serviceCenter/v1"
	utilApi "gitee.com/moyusir/util/api/util/v1"
	"github.com/gorilla/websocket"
	"github.com/imroc/req/v3"
	"google.golang.org/protobuf/types/known/durationpb"
	"net/http"
	"os"
	"testing"
	"time"
)

func TestUserUsecase(t *testing.T) {
	const (
		KONG_HTTP_URL = "http://kong.test.svc.cluster.local:8000"
		username      = "testreg"
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
			Id:       username,
			Password: username,
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
			Id:       username,
			Password: username,
		})
		if err != nil {
			t.Error(err)
		}
	})

	reply, err := userHTTPClient.Login(context.Background(), &utilApi.User{
		Id:       username,
		Password: username,
	})
	if err != nil {
		t.Fatal(err)
	}

	// 等待路由注册生效
	time.Sleep(5 * time.Second)

	// 测试获得客户端文件保存到本地
	t.Run("Test_GetClientCode", func(t *testing.T) {
		client := req.C().DevMode().SetBaseURL("http://localhost:8000")
		response, err := client.R().SetPathParam("username", username).
			Get("/users/client-code/{username}")
		if err != nil {
			t.Fatal(err)
		}
		if response.IsError() {
			t.Fatal(response.Error())
		}

		file, err := os.Create("client_code.zip")
		if err != nil {
			t.Fatal(err)
		}
		defer file.Close()
		_, err = file.Write(response.Bytes())
		if err != nil {
			t.Fatal(err)
		}
	})

	t.Run("Test_WarningPushWebsocket", func(t *testing.T) {
		// 建立接收故障信息推送的ws连接
		conn, _, err := websocket.DefaultDialer.Dial(
			"ws://kong.test.svc.cluster.local:8000/warnings/push/"+username,
			http.Header{
				"X-Api-Key":      {reply.Token},
				"X-Service-Type": {username + "-dp"},
			},
		)
		if err != nil {
			t.Fatal(err)
		}
		err = conn.Close()
		if err != nil {
			t.Error(err)
		}
	})

	// 发送查询设备状态注册信息的请求
	//t.Run("Test_GetDeviceStateRegisterInfo", func(t *testing.T) {
	//	client := req.C().DevMode().SetBaseURL(KONG_HTTP_URL)
	//	response, err := client.R().
	//		SetHeaders(map[string]string{
	//			"X-Api-Key":      reply.Token,
	//			"X-Service-Type": username + "-dp",
	//		}).Get("/register-info/states/0")
	//	if err != nil {
	//		t.Fatal(err)
	//	}
	//	if response.IsError() {
	//		t.Fatal(response.Error())
	//	}
	//
	//	registerInfo := new(utilApi.DeviceStateRegisterInfo)
	//	err = encoding.GetCodec("json").Unmarshal(response.Bytes(), registerInfo)
	//	if err != nil {
	//		t.Fatal(err)
	//	}
	//	// 比较查询得到的注册信息和一开始注册上传的是否一样
	//	if !proto.Equal(registerReq, stateInfo[0]) {
	//		t.Errorf("wrong device state register info %v\n%v", registerReq, stateInfo[0])
	//	}
	//})
}
func TestUser_Register(t *testing.T) {
	// 测试用户利用违规的注册信息进行登录，校验是否成功
	userHTTPClient := StartServiceCenterServer(t)
	username := "testreg"
	t.Cleanup(func() {
		userHTTPClient.Unregister(context.Background(), &utilApi.User{
			Id:       username,
			Password: username,
		})
	})

	t.Run("Test_RepeatedFieldName", func(t *testing.T) {
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
					{
						Name: "power",
						Type: utilApi.Type_DOUBLE,
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
						Name: "time",
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
		}

		// 发送包含重复字段设备状态信息的请求
		repeatStateFieldNameReq := &v1.RegisterRequest{
			User: &utilApi.User{
				Id:       username,
				Password: username,
			},
			DeviceConfigRegisterInfos: configInfo,
			DeviceStateRegisterInfos:  stateInfo,
		}
		_, err := userHTTPClient.Register(context.Background(), repeatStateFieldNameReq)
		if err == nil {
			t.Error("failed to validate repeat field")
		}

		// 发送包含重复字段配置状态信息的请求
		stateInfo[0].Fields[2].Name = "current"
		configInfo[0].Fields[2].Name = "status"
		repeatConfigFieldNameReq := &v1.RegisterRequest{
			User: &utilApi.User{
				Id:       username,
				Password: username,
			},
			DeviceConfigRegisterInfos: configInfo,
			DeviceStateRegisterInfos:  stateInfo,
		}
		_, err = userHTTPClient.Register(context.Background(), repeatConfigFieldNameReq)
		if err == nil {
			t.Error("failed to validate repeat field")
		}
	})

	t.Run("Test_MissingField", func(t *testing.T) {
		// 定义生成代码需要使用的注册信息
		configInfo := []*utilApi.DeviceConfigRegisterInfo{
			{
				Fields: []*utilApi.DeviceConfigRegisterInfo_Field{
					{
						Name: "lost_id",
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
				},
			},
		}

		// 发送不包含id的配置信息的注册请求
		lostConfigIDReq := &v1.RegisterRequest{
			User: &utilApi.User{
				Id:       username,
				Password: username,
			},
			DeviceConfigRegisterInfos: configInfo,
			DeviceStateRegisterInfos:  stateInfo,
		}
		_, err := userHTTPClient.Register(context.Background(), lostConfigIDReq)
		if err == nil {
			t.Error("failed to validate missing id field")
		}

		// 发送不包含id的设备状态信息的注册请求
		configInfo[0].Fields[0].Name = "id"
		stateInfo[0].Fields[0].Name = "lost_id"
		lostStateIDReq := &v1.RegisterRequest{
			User: &utilApi.User{
				Id:       username,
				Password: username,
			},
			DeviceConfigRegisterInfos: configInfo,
			DeviceStateRegisterInfos:  stateInfo,
		}
		_, err = userHTTPClient.Register(context.Background(), lostStateIDReq)
		if err == nil {
			t.Error("failed to validate missing id field")
		}

		// 发送不包含time的设备状态信息的注册请求
		stateInfo[0].Fields[0].Name = "id"
		stateInfo[0].Fields[1].Name = "lost_time"
		lostStateTimeReq := &v1.RegisterRequest{
			User: &utilApi.User{
				Id:       username,
				Password: username,
			},
			DeviceConfigRegisterInfos: configInfo,
			DeviceStateRegisterInfos:  stateInfo,
		}
		_, err = userHTTPClient.Register(context.Background(), lostStateTimeReq)
		if err == nil {
			t.Error("failed to validate missing time field")
		}
	})

	t.Run("Test_WrongWarningRule", func(t *testing.T) {
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
		}
		stateInfo := []*utilApi.DeviceStateRegisterInfo{
			{
				Fields: []*utilApi.DeviceStateRegisterInfo_Field{
					{
						Name: "id",
						Type: utilApi.Type_STRING,
						WarningRule: &utilApi.DeviceStateRegisterInfo_WarningRule{
							CmpRule:              nil,
							AggregationOperation: 0,
							Duration:             nil,
						},
					},
					{
						Name: "time",
						Type: utilApi.Type_TIMESTAMP,
						WarningRule: &utilApi.DeviceStateRegisterInfo_WarningRule{
							CmpRule:              nil,
							AggregationOperation: 0,
							Duration:             nil,
						},
					},
				},
			},
		}

		// 发送为非数值字段注册预警规则的请求
		lostConfigIDReq := &v1.RegisterRequest{
			User: &utilApi.User{
				Id:       username,
				Password: username,
			},
			DeviceConfigRegisterInfos: configInfo,
			DeviceStateRegisterInfos:  stateInfo,
		}
		_, err := userHTTPClient.Register(context.Background(), lostConfigIDReq)
		if err == nil {
			t.Error("failed to validate wrong warning rule")
		}
	})
}
