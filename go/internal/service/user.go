package service

import (
	"context"
	"fancy-todo/internal/config"
	"fancy-todo/internal/repository"
	"fmt"
)

func NewUserService(env *config.Env, userRepo UserRepo) *UserService {
	return &UserService{
		Env: env,
		UserRepo: userRepo,
	}
}

type UserService struct {
	Env *config.Env
	UserRepo UserRepo
}

func (us *UserService) Register(ctx context.Context, data UserServiceRegisterInput) (int, string, error) {
	fmt.Println("User Service: login")
	us.UserRepo.Create(ctx, repository.CreateUserInput{})
	return 0, "", nil
}