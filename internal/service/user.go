package service

import (
	"context"
	"gitee.com/moyusir/service-centre/internal/biz"

	pb "gitee.com/moyusir/service-centre/api/serviceCenter/v1"
	utilApi "gitee.com/moyusir/util/api/util/v1"
)

type UserService struct {
	pb.UnimplementedUserServer
	uc *biz.UserUsecase
}

func NewUserService(uc *biz.UserUsecase) *UserService {
	return &UserService{uc: uc}
}

func (s *UserService) Register(ctx context.Context, req *pb.RegisterRequest) (*pb.RegisterReply, error) {
	token, err := s.uc.Register(req)
	if err != nil {
		return nil, err
	}

	return &pb.RegisterReply{
		Success: true,
		Token:   token,
	}, nil
}

func (s *UserService) GetRegisterInfo(ctx context.Context, req *pb.GetRegisterInfoRequest) (*pb.GetRegisterInfoReply, error) {
	reply := &pb.GetRegisterInfoReply{}
	err := s.uc.GetUserRegisterInfo(req.Token, reply)
	if err != nil {
		return nil, err
	}

	return reply, nil
}

func (s *UserService) Login(ctx context.Context, req *utilApi.User) (*pb.LoginReply, error) {
	token, err := s.uc.Login(req.Id, req.Password)
	if err != nil {
		return nil, err
	}

	return &pb.LoginReply{
		Success: true,
		Token:   token,
	}, nil
}
func (s *UserService) Unregister(ctx context.Context, req *utilApi.User) (*pb.UnregisterReply, error) {
	err := s.uc.Unregister(req.Id, req.Password)
	if err != nil {
		return nil, err
	}

	return &pb.UnregisterReply{
		Success: true,
	}, nil
}

func (s *UserService) DownloadClientCode(ctx context.Context, req *pb.DownloadClientCodeRequest) (*pb.File, error) {
	code, err := s.uc.GetClientCode(req.Username)
	if err != nil {
		return nil, err
	}

	return &pb.File{Content: code, Name: "client_code.zip"}, nil
}
