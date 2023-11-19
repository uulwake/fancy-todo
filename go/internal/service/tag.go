package service

import (
	"context"
	"fancy-todo/internal/config"
	"fancy-todo/internal/repository"
)

func NewTagService(env *config.Env, tagRepo ITagRepo) *TagService {
	return &TagService{
		env: env,
		tagRepo: tagRepo,
	}
}

type TagService struct {
	env *config.Env
	tagRepo ITagRepo
}

func (ts *TagService) Create(ctx context.Context, data TagCreateData) (int64, error) {
	return ts.tagRepo.Create(ctx, repository.TagCreateData{
		Name: data.Name,
		TaskId: data.TaskId,
		UserId: data.UserId,
	})
}

func (ts *TagService) AddExistingTagToTask(ctx context.Context, userId int64, tagId int64, taskId int64) error {
	return ts.tagRepo.AddExistingTagToTask(ctx, userId, tagId, taskId)
}