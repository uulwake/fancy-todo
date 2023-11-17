package service

type TaskServiceCreateInput struct {
	UserId int64
	Title string
	Description string
	TagIDs []int64
}