package repository

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/json"
	"fancy-todo/internal/config"
	"fancy-todo/internal/constant"
	"fancy-todo/internal/database"
	"fancy-todo/internal/libs"
	"fancy-todo/internal/model"
	"time"

	"github.com/huandu/go-sqlbuilder"
	"gopkg.in/guregu/null.v4"
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
	tx, err := tr.db.Pg.Begin()
	if err != nil {
		return 0, libs.DefaultInternalServerError(err)
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

	err = tx.QueryRow(taskQuery, taskArgs...).Scan(&taskId)
	if err != nil {
		return 0, libs.DefaultInternalServerError(err)
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
		_, err = tx.Exec(tagQuery, tagArgs...)
		if err != nil {
			return 0, libs.DefaultInternalServerError(err)
		}
	}


	task := model.Task{
		ID: taskId,
		UserID: data.UserId,
		Title: data.Title,
		Description: &null.String{
			NullString: sql.NullString{String: data.Description, Valid: true},
		},
		Status: constant.TASK_STATUS_ON_GOING,
		CreatedAt: &now,
		UpdatedAt: &now,
	}
	
	taskJson, err := json.Marshal(task)
	if err != nil {
		return 0, libs.DefaultInternalServerError(err)
	}

	_, err = tr.db.Es.Index(constant.ES_INDEX_TASKS, bytes.NewReader(taskJson))
	if err != nil {
		return 0, libs.DefaultInternalServerError(err)
	}

	err = tx.Commit()
	if err != nil {
		return 0, libs.DefaultInternalServerError(err)
	}
	
	return taskId, nil
}

func (tr *TaskRepo) getTagsByTaskId(ctx context.Context, taskIds []int64) (map[int64][]model.Tag, error) {
	sb := sqlbuilder.PostgreSQL.NewSelectBuilder().
		Select("tags.id", "tags.name", "tasks_tags.task_id").
		From("tags").
		JoinWithOption(sqlbuilder.LeftJoin, "tasks_tags", "tasks_tags.tag_id = tags.id")

	taskIdsInterface := make([]any, len(taskIds))
	for i, taskId := range taskIds {
		taskIdsInterface[i] = taskId
	}

	query, args := sb.
		Where(sb.In("tasks_tags.task_id", taskIdsInterface...)).
		OrderBy("tags.id ASC").
		Build()

	rows, err := tr.db.Pg.Query(query, args...)
	if err != nil {
		return nil, libs.DefaultInternalServerError(err)
	}
	defer rows.Close()
	
	mapTagsByTaskId := make(map[int64][]model.Tag)
	for rows.Next() {
		var tag model.Tag
		var task model.Task
		err := rows.Scan(&tag.ID, &tag.Name, &task.ID)
		if err != nil {
			return nil, libs.DefaultInternalServerError(err)
		}
		_, ok := mapTagsByTaskId[task.ID]
		if !ok {
			mapTagsByTaskId[task.ID] = []model.Tag{tag}
		} else {
			mapTagsByTaskId[task.ID] = append(mapTagsByTaskId[task.ID], tag)
		}
	}

	return mapTagsByTaskId, nil
}

func (tr *TaskRepo) GetDetail(ctx context.Context, userId int64, taskId int64) (model.Task, error) {
	sb := sqlbuilder.PostgreSQL.NewSelectBuilder().
		Select("id", "user_id", "title", "description", "status", `"order"`, "created_at", "updated_at").
		From("tasks")

	query, args := sb.
		Where(sb.Equal("id", taskId)).
		Where(sb.Equal("user_id", userId)).
		Build()

	var task model.Task
	err := tr.db.Pg.QueryRow(query, args...).Scan(&task.ID, &task.UserID, &task.Title, &task.Description, &task.Status, &task.Order, &task.CreatedAt, &task.UpdatedAt)
	if err != nil {
		return task, libs.DefaultInternalServerError(err)
	}

	mapTagByTaskId, err := tr.getTagsByTaskId(ctx, []int64{taskId})
	if err != nil {
		return task, err
	}

	if tags, ok := mapTagByTaskId[taskId]; ok {
		task.Tags = &tags
	} else {
		task.Tags = &[]model.Tag{}
	}

	return task, nil
}