package handler

import "fancy-todo/internal/model"

type TaskHandlerCreateBody struct {
	Title string `json:"title" validation:"required,min=1"`
	Description string `json:"description"`
	TagIDs []int64 `json:"tag_ids"`
}

type TaskHandlerCreateResponseData struct {
	Task model.Task `json:"task"`
}

type TaskHandlerCreateResponse struct {
	Data TaskHandlerCreateResponseData `json:"data"`
}
type TaskHandlerGetDetailResponseData struct {
	Task model.Task `json:"task"`
}

type TaskHandlerGetDetailResponse struct {
	Data TaskHandlerGetDetailResponseData `json:"data"`
}