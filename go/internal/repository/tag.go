package repository

import (
	"bytes"
	"context"
	"encoding/json"
	"fancy-todo/internal/config"
	"fancy-todo/internal/constant"
	"fancy-todo/internal/database"
	"fancy-todo/internal/libs"
	"fancy-todo/internal/model"
	"fmt"
	"net/http"
	"time"

	"github.com/huandu/go-sqlbuilder"
)

func NewTagRepo(env *config.Env, db *database.Db) *TagRepo {
	return &TagRepo{
		env: env,
		db: db,
	}
}

type TagRepo struct {
	env *config.Env
	db *database.Db
}

func (tr *TagRepo) Create(ctx context.Context, data TagCreateData) (int64, error) {
	tx, err := tr.db.Pg.Begin()
	if err != nil {
		return 0, libs.DefaultInternalServerError(err)
	}
	defer tx.Rollback()

	var tagId int64
	now := time.Now()

	ib := sqlbuilder.PostgreSQL.NewInsertBuilder()
	query, args := ib.InsertInto("tags").
		Cols("user_id", "name", "created_at", "updated_at").
		Values(data.UserId, data.Name, now, now).
		SQL("RETURNING id").
		Build()

	fmt.Println(query, args)
	err = tx.QueryRow(query, args...).Scan(&tagId)
	if err != nil {
		return 0, libs.DefaultInternalServerError(err)
	}

	if data.TaskId != 0 {
		ib := sqlbuilder.PostgreSQL.NewInsertBuilder()
		query, args := ib.InsertInto("tasks_tags").
			Cols("task_id", "tag_id", "created_at", "updated_at").
			Values(data.TaskId, tagId, now, now).
			Build()

		_, err = tx.Exec(query, args...)
		if err != nil {
			return 0, libs.DefaultInternalServerError(err)
		}
	}
	
	tag := model.Tag{
		ID: tagId,
		UserID: data.UserId,
		Name: data.Name,
		CreatedAt: &now,
		UpdatedAt: &now,
	}

	var tagJson bytes.Buffer
	err = json.NewEncoder(&tagJson).Encode(&tag)
	if err != nil {
		return 0, libs.DefaultInternalServerError(err)
	}

	response, err := tr.db.Es.Index(constant.ES_INDEX_TAGS, &tagJson)
	if err != nil {
		return 0, libs.DefaultInternalServerError(err)
	}
	fmt.Println(response)
	if response.StatusCode != http.StatusCreated {
		return 0, libs.CustomError{
			HTTPCode: http.StatusBadRequest,
			Message: "error inserting data to ElasticSearch",
		}
	}

	err = tx.Commit()
	if err != nil {
		return 0, libs.DefaultInternalServerError(err)
	}

	return tagId, nil
}