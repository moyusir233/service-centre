package data

import (
	"context"
	"gitee.com/moyusir/service-centre/internal/biz"
	"github.com/go-kratos/kratos/v2/errors"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-redis/redis/v8"
)

const (
	// PSWS_KEY 用户密码hash的key
	PSWS_KEY = "passwords"
	// TOKENS_KEY 用户token hash的key
	TOKENS_KEY = "tokens"
)

// RedisRepo redis数据库操作对象，可以理解为dao
type RedisRepo struct {
	client *Data
	logger *log.Helper
}

// NewRedisRepo 实例化redis数据库操作对象
func NewRedisRepo(data *Data, logger log.Logger) biz.UserRepo {
	return &RedisRepo{
		client: data,
		logger: log.NewHelper(logger),
	}
}

// Login 验证用户账号密码，正确时返回用户token
// 用户的密码以用户账号-用户密码键值对的形式存储在hash中，用户的token也同样
func (r *RedisRepo) Login(username, password string) (token string, err error) {
	psw, err := r.client.HGet(context.Background(), PSWS_KEY, username).Result()
	if err != nil {
		return "", err
	} else if psw != password {
		return "", errors.New(
			400, "LOGIN_FAIL", "User's account or password is wrong")
	}

	token, err = r.client.HGet(context.Background(), TOKENS_KEY, username).Result()
	if err != nil {
		return "", err
	}

	return token, nil
}

// Register 用户注册，并保存用户token
func (r *RedisRepo) Register(username, password, token string) error {
	// 保存用户密码以及token，利用事务保证一并执行
	cmders, err := r.client.TxPipelined(context.Background(), func(p redis.Pipeliner) error {
		p.HSetNX(context.Background(), PSWS_KEY, username, password)
		p.HSetNX(context.Background(), TOKENS_KEY, username, token)
		return nil
	})
	if err != nil {
		return err
	}

	// 检查结果
	for _, cmder := range cmders {
		if cmder.Err() != nil {
			return cmder.Err()
		} else if !cmder.(*redis.BoolCmd).Val() {
			return errors.New(400, "REGISTER_FAIL", "User account already exists")
		}
	}

	return nil
}

// UnRegister 注销账户，清除用户相关的所有redis key
func (r *RedisRepo) UnRegister(username string) error {
	// 利用事务保证全部删除完毕
	cmders, err := r.client.TxPipelined(context.Background(), func(p redis.Pipeliner) error {
		// 删除密码和token
		p.HDel(context.Background(), PSWS_KEY, username)
		p.HDel(context.Background(), TOKENS_KEY, username)

		// 获得然后删除和用户相关的键，包括设备配置信息、状态信息、警告信息等
		keys := p.Keys(context.Background(), username+"*").Val()
		p.Del(context.Background(), keys...)

		return nil
	})
	if err != nil {
		return err
	}

	for _, cmder := range cmders {
		if cmder.Err() != nil {
			return cmder.Err()
		}
	}

	return nil
}
