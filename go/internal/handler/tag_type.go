package handler

import "fancy-todo/internal/model"

type TagCreateBody struct {
	Name string	`json:"name" validate:"required"`
	TaskId int64 `json:"task_id"`
}

type TagCreateResponse struct {
	Data TagCreateData `json:"data"`
}

type TagCreateData struct {
	Tag model.Tag `json:"tag"`
}

type TagAddExistingToTaskResponse struct {
	Data TagAddExistingToTaskData `json:"data"`
}

type TagAddExistingToTaskData struct {
	Tag model.Tag `json:"tag"`
}