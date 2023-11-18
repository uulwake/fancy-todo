package handler

import (
	"fancy-todo/internal/config"
	"fancy-todo/internal/model"
	"fancy-todo/internal/service"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
)

func InitUserHandler(echoGroup *echo.Group, env *config.Env, validate *validator.Validate, userService IUserService) {
	uh := &UserHandler{
		echoGroup: echoGroup, 
		env: env, 
		validate: validate, 
		userService: userService,
	}

	uh.echoGroup.POST("/register", uh.Register)
	uh.echoGroup.POST("/login", uh.Login)
}

type UserHandler struct {
	echoGroup *echo.Group
	env *config.Env
	validate *validator.Validate
	userService IUserService
}

func (uh *UserHandler) Register(c echo.Context) error {
	var body UserRegisterRequest
	ctx, err := PreprocessedRequest(c, uh.validate, &body)
	if err != nil {
		return err
	}

	userId, err := uh.userService.Register(ctx, service.UserRegisterInput{Name: body.Name, Email: body.Email, Password: body.Password})
	if err != nil {
		return err
	}

	jwtToken, err := uh.userService.CreateJwtToken(ctx, userId, body.Email)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, UserRegisterLoginResponse{
		Data: UserRegisterLoginResponseData{
			User: model.User{
				ID: userId,
			},
			Token: jwtToken,
		},
	})
}

func (uh *UserHandler) Login(c echo.Context) error {
	var body UserLoginRequest
	ctx, err := PreprocessedRequest(c, uh.validate, &body)
	if err != nil {
		return err
	}

	userId, err := uh.userService.Login(ctx, service.UserLoginInput{
		Email: body.Email,
		Password: body.Password,
	})
	if err != nil {
		return err
	}


	jwtToken, err := uh.userService.CreateJwtToken(ctx, userId, body.Email)
	if err != nil {
		return err
	}


	return c.JSON(http.StatusOK, UserRegisterLoginResponse{
		Data: UserRegisterLoginResponseData{
			User: model.User{
				ID: userId,
			},
			Token: jwtToken,
		},
	})
}