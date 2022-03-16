package biz

import (
	"github.com/go-kratos/kratos/v2/log"
)

type UserUsecase struct {
	repo   UserRepo
	logger *log.Helper
}
type UserRepo interface {
	// Login 用户登录
	Login(username, password string) (token string, err error)
	// Register 用户注册
	Register(username, password, token string) error
	// UnRegister 用户注销
	UnRegister(username string) error
}

func NewUserUsecase(repo UserRepo, logger log.Logger) *UserUsecase {
	return &UserUsecase{
		repo:   repo,
		logger: log.NewHelper(logger),
	}
}
