package handler

import (
	"context"
	"fancy-todo/internal/service"
)

type UserService interface {
	Register(ctx context.Context, data service.UserServiceRegisterInput) (int, string, error)
}