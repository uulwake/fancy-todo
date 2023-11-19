package handler

import (
	"fancy-todo/internal/config"
	"fancy-todo/internal/handler/internal/middleware"
	"fancy-todo/internal/model"
	"fancy-todo/internal/service"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
)

func InitTagHandler(echoGroup *echo.Group, env *config.Env, validate *validator.Validate, tagService ITagService) {
	th := &TagHandler{
		echoGroup: echoGroup,
		env: env,
		validate: validate,
		tagService: tagService,
	}

	th.echoGroup.Use(middleware.AuthenticateJwt(env))
	th.echoGroup.POST("", th.Create)
	th.echoGroup.PATCH("/:tagId/tasks/:taskId", th.AddExistingTagToTask)
	th.echoGroup.GET("/search", th.Search)
	th.echoGroup.DELETE("/:tagId", th.DeleteById)
}

type TagHandler struct {
	echoGroup *echo.Group
	env *config.Env
	validate *validator.Validate
	tagService ITagService
}

func (th *TagHandler) Create(c echo.Context) error {
	var body TagCreateBody
	ctx, err := PreprocessedRequest(c, th.validate, &body)
	if err != nil {
		return err
	}

	userId, err := GetUserIdFromContext(c)
	if err != nil {
		return err
	}

	tagId, err := th.tagService.Create(ctx, service.TagCreateData{
		Name: body.Name,
		UserId: userId,
		TaskId: body.TaskId,
	})
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, TagCreateResponse{
		Data: TagCreateData{
			Tag: model.Tag{
				ID: tagId,
			},
		},
	})
	
}

func (th *TagHandler) AddExistingTagToTask(c echo.Context) error {
	return nil
}

func (th *TagHandler) Search(c echo.Context) error {
	return nil
}

func (th *TagHandler) DeleteById(c echo.Context) error {
	return nil
}

