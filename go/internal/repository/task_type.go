package repository

type TaskCreateInput struct {
	UserId int64
	Title string
	Description string
	TagIDs []int64
}

type TaskGetListsQuery struct {
	Limit int
	Offset int
	SortBy string
	SortOrder string
	Status string
	TagId int64
}
type TaskGetTotalQuery struct {
	Status string
	TagId int64
}

type TaskUpdateByIdInput struct {
	Title string
	Description string
	Status string
	Order int
}