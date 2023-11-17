package service

import (
	"context"
	"fancy-todo/internal/config"
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

func (ts *TaskService) Create(ctx context.Context, data TaskServiceCreateInput) (int64, error) {
	return ts.taskRepo.Create(ctx, repository.TaskRepoCreateInput{
		UserId: data.UserId,
		Title: data.Title,
		Description: data.Description,
		TagIDs: data.TagIDs,
	})
}