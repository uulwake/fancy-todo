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

func NewTaskRepo(env *config.Env, db *database.Db) *TaskRepo {
	return &TaskRepo{
		env: env,
		db: db,
	}
}

type TaskRepo struct {
	env *config.Env
	db *database.Db
}

func (tr *TaskRepo) Create(ctx context.Context, data TaskRepoCreateInput) (int64, error) {
	fmt.Println("TaskRepo::Create")
	fail := func(err error) (int64, error) {
		return 0, libs.CustomError{
			HTTPCode: http.StatusInternalServerError,
			Message: err.Error(),
		}
	}
	
	tx, err := tr.db.Pg.Begin()
	if err != nil {
		return fail(err)
	}
	defer tx.Rollback()

	var taskId int64
	now := time.Now()
	
	taskQuery, taskArgs := sqlbuilder.PostgreSQL.
		NewInsertBuilder().
		InsertInto("tasks").
		Cols("user_id", "title", "description", "status", "created_at", "updated_at").
		Values(data.UserId, data.Title, data.Description,  constant.TASK_STATUS_ON_GOING,now, now).
		SQL("RETURNING id").
		Build()

	fmt.Println(taskQuery)
	
	err = tx.QueryRow(taskQuery, taskArgs...).Scan(&taskId)
	if err != nil {
		return fail(err)
	}

	if len(data.TagIDs) > 0 {
		sb := sqlbuilder.PostgreSQL.
			NewInsertBuilder().
			InsertInto("tasks_tags").
			Cols("task_id", "tag_id", "created_at", "updated_at")

		for _, tagId := range data.TagIDs {
			sb.Values(taskId, tagId, now, now)
		}

		tagQuery, tagArgs := sb.Build()
		fmt.Println(tagQuery)
		_, err = tx.Exec(tagQuery, tagArgs...)
		if err != nil {
			return fail(err)
		}
	}


	task := model.Task{
		ID: taskId,
		UserID: data.UserId,
		Title: data.Title,
		Description: data.Description,
		Status: constant.TASK_STATUS_ON_GOING,
		CreatedAt: now,
		UpdatedAt: now,
	}
	
	taskJson, err := json.Marshal(task)
	if err != nil {
		return fail(err)
	}

	_, err = tr.db.Es.Index(constant.ES_INDEX_TASKS, bytes.NewReader(taskJson))
	if err != nil {
		return fail(err)
	}

	err = tx.Commit()
	if err != nil {
		return fail(err)
	}
	
	return taskId, nil
}