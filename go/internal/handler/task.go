package handler

import (
	"context"
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
	var body TaskCreateBody
	ctx, err := PreprocessedRequest(c, th.validate, &body)
	if err != nil {
		return err
	}

	userId, err := GetUserIdFromContext(c)
	if err != nil {
		return err
	}
	
	taskId, err := th.taskService.Create(ctx, service.TaskCreateInput{
		UserId: userId,
		Title: body.Title,
		Description: body.Description,
		TagIDs: body.TagIDs,
	})

	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, TaskCreateResponse{
		Data: TaskCreateResponseData{
			Task: model.Task{
				ID: taskId,
			},
		},
	})
}

func (th *TaskHandler) GetLists(c echo.Context) error {
	ctx, err := PreprocessedRequest(c, th.validate, nil)
	if err != nil {
		return err
	}

	userId, err := GetUserIdFromContext(c)
	if err != nil {
		return err
	}

	commonQueryParam, err := ConvertCommonQueryParam(c)
	if err != nil {
		return err
	}

	status := c.QueryParam("status")
	tagIdStr := c.QueryParam("tag_id")

	var tagId int64
	if tagIdStr != "" {
		tagId, err = strconv.ParseInt(tagIdStr, 10, 64)
		if err != nil {
			return libs.CustomError{
				HTTPCode: http.StatusBadRequest,
				Message: fmt.Sprintf("invalid tag ID %s", tagIdStr),
			}
		}
	} 

	queryParamLists := TaskGetListsQueryParam{
		QueryParam: commonQueryParam,
		Status: status,
		TagId: tagId,
	}

	queryParamTotal := TaskGetTotalQueryParam{
		Status: status,
		TagId: tagId,
	}

	getListsOutputChan := make(chan []model.Task, 1)
	getTotalOutputChan := make(chan int64, 1)
	errorChan := make(chan error, 2)

		
	go func(getListsOutputChan chan<- []model.Task, errorChan chan<- error, ctx context.Context, userId int64, queryParam TaskGetListsQueryParam) {
		tasks, err := th.taskService.GetLists(ctx, userId, service.TaskGetListsQueryParam{
			QueryParam: queryParam.QueryParam,
			Status: queryParam.Status,
			TagId: queryParam.TagId,
		})

		getListsOutputChan <- tasks
		errorChan <- err
		
		close(getListsOutputChan)
	}(getListsOutputChan, errorChan, ctx, userId, queryParamLists)


	go func(getTotalOutputChan chan<- int64, errorChan chan<- error, ctx context.Context, userId int64, queryParam TaskGetTotalQueryParam) {
		total, err := th.taskService.GetTotal(ctx, userId, service.TaskGetTotalQueryParam{
			Status: queryParam.Status,
			TagId: queryParam.TagId,
		})

		getTotalOutputChan <- total
		errorChan <- err

		close(getTotalOutputChan)
	}(getTotalOutputChan, errorChan, ctx, userId, queryParamTotal)


	tasks := <-getListsOutputChan
	total := <-getTotalOutputChan

	err1 := <-errorChan
	err2 := <-errorChan 
	close(errorChan)

	if err1 != nil {
		return err1
	}

	if err2 != nil {
		return err2
	}


	return c.JSON(http.StatusOK, TaskGetListsResponse{
		Data: TaskGetListData{
			Tasks: tasks,
		},
		Page: libs.Pagination{
			Size: queryParamLists.PageSize,
			Number: queryParamLists.PageNumber,
			Total: total,
		},
	})	
}

func (th *TaskHandler) Search(c echo.Context) error {
	ctx, err := PreprocessedRequest(c, th.validate, nil)
	if err != nil {
		return err
	}

	userId, err := GetUserIdFromContext(c)
	if err != nil {
		return err
	}

	title := c.QueryParam("title")

	tasks, err := th.taskService.Search(ctx, userId, title)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, TaskSearchResponse{
		Data: TaskSearchData{
			Tasks: tasks,
		},
	})
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
			Message: fmt.Sprintf("invalid task ID %s", taskIdParam),
		}
	}

	task, err := th.taskService.GetDetail(ctx, userId, taskId)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, TaskGetDetailResponse{
		Data: TaskGetDetailResponseData{
			Task: task,
		},
	})
}

func (th *TaskHandler) UpdateById(c echo.Context) error {
	var body TaskUpdateByIdBody
	ctx, err := PreprocessedRequest(c, th.validate, &body)
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
			Message: fmt.Sprintf("invalid task ID %s", taskIdParam),
		}
	}

	err = th.taskService.UpdateById(ctx, userId, taskId, service.TaskUpdateByIdInput{
		Title: body.Title,
		Description: body.Description,
		Status: body.Status,
		Order: body.Order,
	})
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, TaskUpdateByIdResponse{
		Data: TaskUpdateByIdData{
			Task: model.Task{
				ID: taskId,
			},
		},
	})
}

func (th *TaskHandler) DeleteById(c echo.Context) error {
	return nil
}