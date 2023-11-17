package repository

type TaskRepoCreateInput struct {
	UserId int64
	Title string
	Description string
	TagIDs []int64
}