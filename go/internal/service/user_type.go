package service

type UserServiceRegisterInput struct {
	Name string
	Email string
	Password string
}

type UserServiceLoginInput struct {
	Email string
	Password string
}