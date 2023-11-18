package repository

type UserCreateInput struct {
	Name string
	Email string
	Password string
}

type UserGetDetailInput struct {
	ID int64
	Email string
	Cols []string
	Values []any
}
