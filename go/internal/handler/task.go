package handler

import (
	"fancy-todo/internal/config"
	"fancy-todo/internal/handler/internal/middleware"
	"fancy-todo/internal/model"
	"fancy-todo/internal/service"
	"fmt"
	"net/http"
	"reflect"

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
	fmt.Println("TaskHandler::Create")
	var body TaskHandlerCreateBody
	ctx, err := PreprocessedRequest(c, th.validate, &body)
	if err != nil {
		return err
	}
	
	fmt.Println(c.Get("user_id"), reflect.TypeOf(c.Get("user_id")).Kind())
	taskId, err := th.taskService.Create(ctx, service.TaskServiceCreateInput{
		UserId: c.Get("user_id").(int64),
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
	return nil
}

func (th *TaskHandler) UpdateById(c echo.Context) error {
	return nil
}

func (th *TaskHandler) DeleteById(c echo.Context) error {
	return nil
}