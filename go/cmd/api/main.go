package main

import (
	"fancy-todo/internal/config"
	"fancy-todo/internal/database"
	"fancy-todo/internal/handler"
	"fancy-todo/internal/repository"
	"fancy-todo/internal/service"
	"log"

	"github.com/labstack/echo/v4"
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

	// service
	userService := service.NewUserService(env, userRepo)

	// handler
	e := echo.New()

	v1Group := e.Group("/v1")
	handler.InitUserHandler(v1Group.Group("/users"), env, userService)

	e.Logger.Fatal(e.Start(":3001"))
}