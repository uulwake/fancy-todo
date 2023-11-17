package repository

import (
	"context"
	"fancy-todo/internal/config"
	"fancy-todo/internal/database"
	"fancy-todo/internal/libs"
	"net/http"
	"time"

	"github.com/huandu/go-sqlbuilder"
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

func (ur *UserRepo) Create(ctx context.Context, data CreateUserInput) (int64, error) {
	now := time.Now()
	sb := sqlbuilder.PostgreSQL.NewInsertBuilder()
	sb.InsertInto("users").Cols("name", "email", "password", "created_at", "updated_at")
	sb.Values(data.Name, data.Email, data.Password, now, now)
	sb.SQL("RETURNING id")

	query, args := sb.Build()

	var userId int64
	err := ur.db.Pg.QueryRow(query, args...).Scan(&userId)
	if err != nil {
		return 0, libs.CustomError{
			HTTPCode: http.StatusInternalServerError,
			Message: err.Error(),
		}
	}

	return userId, nil
}
