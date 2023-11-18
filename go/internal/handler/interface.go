package handler

import (
	"context"
	"fancy-todo/internal/model"
	"fancy-todo/internal/service"
)

type IUserService interface {
	CreateJwtToken(ctx context.Context, userId int64, email string) (string, error)
	Register(ctx context.Context, data service.UserServiceRegisterInput) (int64, error)
	Login(ctx context.Context, data service.UserServiceLoginInput) (int64, error)
}

type ITaskService interface {
	Create(ctx context.Context, data service.TaskServiceCreateInput) (int64, error)
	GetDetail(ctx context.Context, userId int64, taskId int64) (model.Task, error)
}