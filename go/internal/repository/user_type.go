package repository

type CreateUserInput struct {
	Name string
	Email string
	Password string
}

type GetDetailUserInput struct {
	ID int64
	Email string
	Cols []string
	Values []any
}
