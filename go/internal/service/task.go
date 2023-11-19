package service

import (
	"context"
	"fancy-todo/internal/config"
	"fancy-todo/internal/model"
	"fancy-todo/internal/repository"
)

func NewTaskService(env *config.Env, taskRepo ITaskRepo) *TaskService {
	return &TaskService{
		env: env,
		taskRepo: taskRepo,
	}
}

type TaskService struct {
	env *config.Env
	taskRepo ITaskRepo
}

func (ts *TaskService) Create(ctx context.Context, data TaskCreateInput) (int64, error) {
	return ts.taskRepo.Create(ctx, repository.TaskCreateInput{
		UserId: data.UserId,
		Title: data.Title,
		Description: data.Description,
		TagIDs: data.TagIDs,
	})
}

func (ts *TaskService) GetDetail(ctx context.Context, userId int64, taskId int64) (model.Task, error) {
	return ts.taskRepo.GetDetail(ctx, userId, taskId)
}

func (ts *TaskService) GetLists(ctx context.Context, userId int64, queryParam TaskGetListsQueryParam) ([]model.Task, error) {
	return ts.taskRepo.GetLists(ctx, userId, repository.TaskGetListsQuery{
		Limit: queryParam.PageSize,
		Offset: (queryParam.PageNumber - 1) * queryParam.PageSize,
		SortBy: queryParam.SortKey,
		SortOrder: queryParam.SortOrder,
		Status: queryParam.Status,
		TagId: queryParam.TagId,
	})
}

func (ts *TaskService) GetTotal(ctx context.Context, userId int64, queryParam TaskGetTotalQueryParam) (int64, error) {
	return ts.taskRepo.GetTotal(ctx, userId, repository.TaskGetTotalQuery{
		Status: queryParam.Status,
		TagId: queryParam.TagId,
	})
}

func (ts *TaskService) Search(ctx context.Context, userId int64, title string) ([]model.Task, error) {
	return ts.taskRepo.Search(ctx, userId, title)
}