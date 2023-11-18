package handler

import (
	"fancy-todo/internal/config"
	"fancy-todo/internal/handler/internal/middleware"
	"fancy-todo/internal/libs"
	"fancy-todo/internal/model"
	"fancy-todo/internal/service"
	"fmt"
	"net/http"
	"strconv"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
)

func InitTaskHandler(echoGroup *echo.Group, env *config.Env, validate *validator.Validate, taskService ITaskService) {
	th := &TaskHandler{
		echoGroup: echoGroup,
		env: env,
		validate: validate,
		taskService: taskService,
	}

	th.echoGroup.Use(middleware.AuthenticateJwt(env))
	th.echoGroup.POST("", th.Create)
	th.echoGroup.GET("", th.GetLists)
	th.echoGroup.GET("/:taskId", th.GetDetailById)
	th.echoGroup.GET("/search", th.Search)
	th.echoGroup.PATCH("/:taskId", th.UpdateById)
	th.echoGroup.DELETE("/:taskId", th.DeleteById)
}

type TaskHandler struct {
	echoGroup *echo.Group
	env *config.Env
	validate *validator.Validate
	taskService ITaskService
}

func (th *TaskHandler) Create(c echo.Context) error {
	var body TaskHandlerCreateBody
	ctx, err := PreprocessedRequest(c, th.validate, &body)
	if err != nil {
		return err
	}

	userId, err := GetUserIdFromContext(c)
	if err != nil {
		return err
	}
	
	taskId, err := th.taskService.Create(ctx, service.TaskServiceCreateInput{
		UserId: userId,
		Title: body.Title,
		Description: body.Description,
		TagIDs: body.TagIDs,
	})

	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, TaskHandlerCreateResponse{
		Data: TaskHandlerCreateResponseData{
			Task: model.Task{
				ID: taskId,
			},
		},
	})
}

func (th *TaskHandler) GetLists(c echo.Context) error {
	return nil
}

func (th *TaskHandler) Search(c echo.Context) error {
	return nil
}

func (th *TaskHandler) GetDetailById(c echo.Context) error {
	ctx, err := PreprocessedRequest(c, th.validate, nil)
	if err != nil {
		return err
	}

	userId, err := GetUserIdFromContext(c)
	if err != nil {
		return err
	}
	
	taskIdParam := c.Param("taskId")
	taskId, err := strconv.ParseInt(taskIdParam, 10, 64)
	if err != nil {
		return libs.CustomError{
			HTTPCode: http.StatusBadRequest,
			Message: fmt.Sprintf("Invalid task ID %s", taskIdParam),
		}
	}

	task, err := th.taskService.GetDetail(ctx, userId, taskId)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, TaskHandlerGetDetailResponse{
		Data: TaskHandlerGetDetailResponseData{
			Task: task,
		},
	})
}

func (th *TaskHandler) UpdateById(c echo.Context) error {
	return nil
}

func (th *TaskHandler) DeleteById(c echo.Context) error {
	return nil
}