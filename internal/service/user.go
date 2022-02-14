package service

import (
	"context"
	"gitee.com/moyusir/service-centre/internal/biz"

	pb "gitee.com/moyusir/service-centre/api/serviceCenter/v1"
)

type UserService struct {
	pb.UnimplementedUserServer
	uc            *biz.UserUsecase
	codeGenerator *biz.CodeGenerateUsecase
	kubectl       *biz.KubeControlUsecase
}

func NewUserService() *UserService {
	return &UserService{}
}

func (s *UserService) Register(ctx context.Context, req *pb.RegisterRequest) (*pb.RegisterReply, error) {
	return &pb.RegisterReply{}, nil
}
func (s *UserService) Login(ctx context.Context, req *pb.LoginRequest) (*pb.LoginReply, error) {
	return &pb.LoginReply{}, nil
}
