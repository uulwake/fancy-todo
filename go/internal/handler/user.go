package handler

import (
	"encoding/json"
	"fancy-todo/internal/config"
	"fancy-todo/internal/libs"
	"fancy-todo/internal/model"
	"fancy-todo/internal/service"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
)

func InitUserHandler(echoGroup *echo.Group, env *config.Env, validate *validator.Validate, userService UserService) {
	uh := &UserHandler{echoGroup: echoGroup, env: env, validate: validate, userService: userService}
	uh.echoGroup.POST("/register", uh.Register)
}

type UserHandler struct {
	echoGroup *echo.Group
	env *config.Env
	validate *validator.Validate
	userService UserService
}

func (uh *UserHandler) Register(c echo.Context) error {
	ctx := CreateContext(c)
	
	var body UserRegisterRequest
	err := json.NewDecoder(c.Request().Body).Decode(&body)
	if err != nil {
		return libs.CustomError{
			HTTPCode: http.StatusInternalServerError,
			Message: err.Error(),
		}
	}
	
	err = uh.validate.Struct(body)
	if err != nil {
		return libs.CustomError{
			HTTPCode: http.StatusBadRequest,
			Message: err.Error(),
		}
	}

	userId, err := uh.userService.Register(ctx, service.UserServiceRegisterInput{Name: body.Name, Email: body.Email, Password: body.Password})
	if err != nil {
		return err
	}

	jwtToken, err := uh.userService.CreateJwtToken(ctx, userId, body.Email)
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
