package biz

type UserUsecase struct {
	repo UserRepo
}
type UserRepo interface {
}

func NewUserUsecase(repo UserRepo) *UserUsecase {
	return nil
}
