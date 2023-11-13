package service

import (
	"context"
	"fancy-todo/internal/repository"
)

type UserRepo interface {
	Create(ctx context.Context, data repository.CreateUserInput) (int, error)
}