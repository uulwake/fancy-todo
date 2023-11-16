package handler

import "fancy-todo/internal/model"

type UserRegisterRequest struct {
	Name string `json:"name"`
	Email string `json:"email"`
	Password string `json:"password"`
}

type UserRegisterResponseData struct {
	model.User `json:"user"`
	Token string `json:"jwt_token"`
}

type UserRegisterResponse struct {
	Data UserRegisterResponseData `json:"data"`
}
