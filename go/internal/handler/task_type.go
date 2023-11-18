package handler

import (
	"fancy-todo/internal/libs"
	"fancy-todo/internal/model"
)

type TaskCreateBody struct {
	Title string `json:"title" validation:"required,min=1"`
	Description string `json:"description"`
	TagIDs []int64 `json:"tag_ids"`
}

type TaskCreateResponseData struct {
	Task model.Task `json:"task"`
}

type TaskCreateResponse struct {
	Data TaskCreateResponseData `json:"data"`
}
type TaskGetDetailResponseData struct {
	Task model.Task `json:"task"`
}

type TaskGetDetailResponse struct {
	Data TaskGetDetailResponseData `json:"data"`
}

type TaskGetListsQueryParam struct {
	libs.QueryParam
	Status string 
	TagId int64
}
type TaskGetTotalQueryParam struct {
	Status string 
	TagId int64
}

type TaskGetListData struct {
	Tasks []model.Task `json:"tasks"`
}
type TaskGetListsResponse struct {
	Data TaskGetListData `json:"data"`
	Page libs.Pagination `json:"page"`
}