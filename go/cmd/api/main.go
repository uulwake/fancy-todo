package main

import (
	"fancy-todo/internal/config"
	"fancy-todo/internal/database"
	"fancy-todo/internal/handler"
	"fancy-todo/internal/repository"
	"fancy-todo/internal/service"
	"log"
	"net/http"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	env, err := config.NewEnv()
	if err != nil {
		log.Fatal("Failed reading environment variables", err)
	}

	db, err := database.NewDb(env)
	if err != nil {
		log.Fatal("Failed connect to database", err)
	}

	// repository
	userRepo := repository.NewUserRepo(env, db)
	taskRepo := repository.NewTaskRepo(env, db)
	tagRepo := repository.NewTagRepo(env, db)

	// service
	userService := service.NewUserService(env, userRepo)
	taskService := service.NewTaskService(env, taskRepo)
	tagService := service.NewTagService(env, tagRepo)

	// handler
	e := echo.New()
	e.GET("/hc", func(c echo.Context) error {
		return c.JSON(http.StatusOK, struct{ Status string `json:"status"`}{Status: "ok"})
	})

	e.Use(middleware.RequestID())
	e.Use(middleware.TimeoutWithConfig(middleware.TimeoutConfig{
		ErrorMessage: "request timeout",
		Timeout: 30 * time.Second,
	}))
	e.IPExtractor = echo.ExtractIPFromRealIPHeader()
	e.HTTPErrorHandler = handler.CustomHTTPErrorHandler

	validate := validator.New()

	v1Group := e.Group("/v1")
	handler.InitUserHandler(v1Group.Group("/users"), env, validate, userService)
	handler.InitTaskHandler(v1Group.Group("/tasks"), env, validate, taskService)
	handler.InitTagHandler(v1Group.Group("/tags"), env, validate, tagService)

	// run
	e.Logger.Fatal(e.Start(":3001"))
}