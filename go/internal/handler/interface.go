package handler

import (
	"context"
	"fancy-todo/internal/service"
)

type UserService interface {
	CreateJwtToken(ctx context.Context, userId int64, email string) (string, error)
	Register(ctx context.Context, data service.UserServiceRegisterInput) (int64, error)
	Login(ctx context.Context, data service.UserServiceLoginInput) (int64, error)
}