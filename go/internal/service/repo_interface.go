package service

import (
	"context"
	"fancy-todo/internal/model"
	"fancy-todo/internal/repository"
)

type IUserRepo interface {
	Create(ctx context.Context, data repository.UserCreateInput) (int64, error)
	GetDetail(ctx context.Context, queryOption repository.UserGetDetailInput) error
}

type ITaskRepo interface {
	Create(ctx context.Context, data repository.TaskCreateInput) (int64, error)
	GetDetail(ctx context.Context, userId int64, taskId int64) (model.Task, error)
	GetLists(ctx context.Context, userId int64, query repository.TaskGetListsQuery) ([]model.Task, error)
	GetTotal(ctx context.Context, userId int64, query repository.TaskGetTotalQuery) (int64, error)
	Search(ctx context.Context, userId int64, title string) ([]model.Task, error)
	UpdateById(ctx context.Context, userId int64, taskId int64, task repository.TaskUpdateByIdInput) error
	DeleteById(ctx context.Context, userId int64, taskId int64) error 
}

type ITagRepo interface {
	Create(ctx context.Context, data repository.TagCreateData) (int64, error)
	AddExistingTagToTask(ctx context.Context, userId int64, tagId int64, taskId int64) error
}