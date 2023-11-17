package handler

import (
	"context"
	"fancy-todo/internal/service"
)

type UserService interface {
	Register(ctx context.Context, data service.UserServiceRegisterInput) (int64, error)
	CreateJwtToken(ctx context.Context, userId int64, email string) (string, error)
}