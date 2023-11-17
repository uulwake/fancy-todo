package handler

import (
	"encoding/json"
	"fancy-todo/internal/config"
	"fancy-todo/internal/model"
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
	ctx := GenerateContext(c)

	var body UserRegisterRequest
	err := json.NewDecoder(c.Request().Body).Decode(&body)
	if err != nil {
		return err
	}
	
	fmt.Println("UserHandler: Register", body)
	userId, err := uh.UserService.Register(ctx, service.UserServiceRegisterInput{Name: body.Name, Email: body.Email, Password: body.Password})
	if err != nil {
		return err
	}

	jwtToken, err := uh.UserService.CreateJwtToken(ctx, userId, body.Email)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, UserRegisterResponse{
		Data: UserRegisterResponseData{
			User: model.User{
				ID: userId,
			},
			Token: jwtToken,
		},
	})
}
