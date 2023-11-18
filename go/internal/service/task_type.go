package service

import "fancy-todo/internal/libs"

type TaskCreateInput struct {
	UserId int64
	Title string
	Description string
	TagIDs []int64
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