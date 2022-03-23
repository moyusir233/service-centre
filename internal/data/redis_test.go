package data

import (
	"gitee.com/moyusir/service-centre/internal/conf"
	"github.com/go-kratos/kratos/v2/log"
	"testing"
)

func TestRedisRepo(t *testing.T) {
	bootstrap, err := conf.LoadConfig("../../configs/config.yaml")
	if err != nil {
		t.Fatal(err)
	}

	data, cleanUp, err := NewData(bootstrap.Data, log.DefaultLogger)
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(cleanUp)

	redisRepo := NewRedisRepo(data)
	var (
		username = "test"
		password = "test"
		token    = "test"
	)

	// 测试注册
	t.Run("Register", func(t *testing.T) {
		err := redisRepo.Register(username, password, token)
		if err != nil {
			t.Fatal(err)
		}

		// 测试利用相同账号重复注册
		err = redisRepo.Register(username, password, token)
		if err == nil {
			t.Fatal("允许了相同的账号注册")
		}
	})

	// 测试登录
	t.Run("Login", func(t *testing.T) {
		newToken, err := redisRepo.Login(username, password)
		if err != nil {
			t.Fatal(err)
		}

		if newToken != token {
			t.Fatal("登录获得的token与注册时的不一致")
		}
	})

	// 测试注销
	t.Run("Unregister", func(t *testing.T) {
		err := redisRepo.UnRegister(username)
		if err != nil {
			t.Fatal(err)
		}

		// 注销后尝试再次登录
		_, err = redisRepo.Login(username, password)
		if err == nil {
			t.Fatal("注销功能失效")
		}
	})
}
