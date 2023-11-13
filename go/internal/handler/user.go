package handler

import (
	"context"
	"fancy-todo/internal/config"
	"fancy-todo/internal/service"
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
)

func InitUserHandler(echoGroup *echo.Group, env *config.Env, userService UserService) {
	uh := &UserHandler{echoGroup: echoGroup, env: env, UserService: userService}
	uh.echoGroup.POST("/register", uh.Register)
}

type UserHandler struct {
	echoGroup *echo.Group
	env *config.Env
	UserService UserService
}

func (uh *UserHandler) Register(c echo.Context) error {
	fmt.Println("UserHandler: Register")
	uh.UserService.Register(context.TODO(), service.UserServiceRegisterInput{})

	return c.JSON(http.StatusOK, UserRegisterResponse{Message: "success register"})
}
