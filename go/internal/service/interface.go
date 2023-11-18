package service

import (
	"context"
	"fancy-todo/internal/model"
	"fancy-todo/internal/repository"
)

type IUserRepo interface {
	Create(ctx context.Context, data repository.CreateUserInput) (int64, error)
	GetDetail(ctx context.Context, queryOption repository.GetDetailUserInput) error
}

type ITaskRepo interface {
	Create(ctx context.Context, data repository.TaskRepoCreateInput) (int64, error)
	GetDetail(ctx context.Context, userId int64, taskId int64) (model.Task, error)
}