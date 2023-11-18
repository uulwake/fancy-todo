package repository

import (
	"context"
	"fancy-todo/internal/config"
	"fancy-todo/internal/database"
	"fancy-todo/internal/libs"
	"time"

	"github.com/huandu/go-sqlbuilder"
)

func NewUserRepo(env *config.Env, db *database.Db) *UserRepo {
	return &UserRepo{
		env: env,
		db: db,
	}
}

type UserRepo struct {
	env *config.Env
	db *database.Db
}

func (ur *UserRepo) Create(ctx context.Context, data CreateUserInput) (int64, error) {
	now := time.Now()

	query, args := sqlbuilder.PostgreSQL.NewInsertBuilder().
		InsertInto("users").
		Cols("name", "email", "password", "created_at", "updated_at").
		Values(data.Name, data.Email, data.Password, now, now).
		SQL("RETURNING id").
		Build()

	var userId int64
	err := ur.db.Pg.QueryRow(query, args...).Scan(&userId)
	if err != nil {
		return 0, libs.DefaultInternalServerError(err)
	}

	return userId, nil
}

func (ur *UserRepo) GetDetail(ctx context.Context, queryOption GetDetailUserInput) error {
	sb := sqlbuilder.PostgreSQL.NewSelectBuilder()
	sb.Select(queryOption.Cols...).From("users").Limit(1)

	if queryOption.ID != 0 {
		sb.Where(sb.Equal("id", queryOption.ID))
	}

	if queryOption.Email != "" {
		sb.Where(sb.Equal("email", queryOption.Email))
	}

	query, args := sb.Build()
	
	err := ur.db.Pg.QueryRow(query, args...).Scan(queryOption.Values...)
	if err != nil {
		return libs.DefaultInternalServerError(err)
	}

	return nil
}
