package repository

type CreateUserInput struct {
	Name string
	Email string
	Password string
}

type GetDetailUserInput struct {
	ID int
	Email string
	Cols []string
}