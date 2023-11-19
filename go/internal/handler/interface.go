package handler

import (
	"context"
	"fancy-todo/internal/model"
	"fancy-todo/internal/service"
)

type IUserService interface {
	CreateJwtToken(ctx context.Context, userId int64, email string) (string, error)
	Register(ctx context.Context, data service.UserRegisterInput) (int64, error)
	Login(ctx context.Context, data service.UserLoginInput) (int64, error)
}

type ITaskService interface {
	Create(ctx context.Context, data service.TaskCreateInput) (int64, error)
	GetDetail(ctx context.Context, userId int64, taskId int64) (model.Task, error)
	GetLists(ctx context.Context, userId int64, queryParam service.TaskGetListsQueryParam) ([]model.Task, error)
	GetTotal(ctx context.Context, userId int64, queryParam service.TaskGetTotalQueryParam) (int64, error)
	Search(ctx context.Context, userId int64, title string) ([]model.Task, error)
	UpdateById(ctx context.Context, userId int64, taskId int64, task service.TaskUpdateByIdInput) error
}