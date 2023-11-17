package repository

import (
	"context"
	"fancy-todo/internal/config"
	"fancy-todo/internal/database"
	"fmt"
)

func NewUserRepo(env *config.Env, db *database.Db) *UserRepo {
	userRepo := &UserRepo{
		env: env,
		db: db,
	}

	return userRepo
}

type UserRepo struct {
	env *config.Env
	db *database.Db
}

func (ur *UserRepo) Create(ctx context.Context, data CreateUserInput) (int, error) {
	fmt.Println("UserRepo: Create")
	return 0, nil
}
