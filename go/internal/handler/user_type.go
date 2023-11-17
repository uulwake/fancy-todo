package handler

import "fancy-todo/internal/model"

type UserRegisterRequest struct {
	Name string `json:"name" validate:"required"`
	Email string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=3,max=20"`
}

type UserRegisterResponseData struct {
	model.User `json:"user"`
	Token string `json:"jwt_token"`
}

type UserRegisterResponse struct {
	Data UserRegisterResponseData `json:"data"`
}
