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
	"fmt"
	"strings"
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

func (tr *TaskRepo) Create(ctx context.Context, data TaskCreateInput) (int64, error) {
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
	mapTagsByTaskId := make(map[int64][]model.Tag)
	if len(taskIds) == 0 {
		return mapTagsByTaskId, nil
	}

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

func (tr *TaskRepo) GetLists(ctx context.Context, userId int64, queryParam TaskGetListsQuery) ([]model.Task, error) {
	sb := sqlbuilder.PostgreSQL.NewSelectBuilder()
		
	sb.Select(
			"tasks.id",
			"tasks.user_id",
			"tasks.title",
			"tasks.status",
			"tasks.order",
			"tasks.created_at",
			"tasks.updated_at",
		).
		From("tasks").
		Where(sb.Equal("tasks.user_id", userId)).
		Limit(queryParam.Limit).
		Offset(queryParam.Offset)

	if queryParam.Status != "" {
		sb.Where(sb.Equal("tasks.status", queryParam.Status))
	}

	if queryParam.SortBy != "" && queryParam.SortOrder != "" && libs.CheckValidSortOrder(queryParam.SortOrder) {
		sb.OrderBy(fmt.Sprintf("%s %s", queryParam.SortBy, strings.ToUpper(queryParam.SortOrder)))
	}

	if (queryParam.TagId != 0) {
		sb.JoinWithOption(sqlbuilder.LeftJoin, "tasks_tags", "tasks_tags.task_id = tasks.id")
		sb.Where(sb.Equal("tasks_tags.tag_id", queryParam.TagId))
	}

	sb.OrderBy("tasks.order ASC NULLS LAST")
	sb.OrderBy("tasks.id ASC")
	
	tasks := []model.Task{}
	query, args := sb.Build()
	rows, err := tr.db.Pg.Query(query, args...)
	if err != nil {
		return tasks, libs.DefaultInternalServerError(err)
	}
	defer rows.Close()

	for rows.Next() {
		var task model.Task
		err := rows.Scan(&task.ID, &task.UserID, &task.Title, &task.Status, &task.Order, &task.CreatedAt, &task.UpdatedAt)
		if err != nil {
			return tasks, libs.DefaultInternalServerError(err)
		}

		tasks = append(tasks, task)
	}


	taskIds := make([]int64, len(tasks))
	for _, task := range tasks {
		taskIds = append(taskIds, task.ID)
	}

	mapTagsByTaskId, err := tr.getTagsByTaskId(ctx, taskIds)
	if err != nil {
		return tasks, libs.DefaultInternalServerError(err)
	}

	for i, task := range tasks {
		if tags, ok := mapTagsByTaskId[task.ID]; ok {
			tasks[i].Tags = &tags
		} else {
			tasks[i].Tags = &[]model.Tag{}
		}
	}

	return tasks, nil
}

func (tr *TaskRepo) GetTotal(ctx context.Context, userId int64, queryParam TaskGetTotalQuery) (int64, error) {
	sb := sqlbuilder.PostgreSQL.NewSelectBuilder()

	sb.Select("COUNT(tasks.id)").
		From("tasks").
		Where(sb.Equal("tasks.user_id", userId))

	if queryParam.Status != "" {
		sb.Where(sb.Equal("tasks.status", queryParam.Status))
	}

	if (queryParam.TagId != 0) {
		sb.JoinWithOption(sqlbuilder.LeftJoin, "tasks_tags", "tasks_tags.task_id = tasks.id")
		sb.Where(sb.Equal("tasks_tags.tag_id", queryParam.TagId))
	}


	var total int64
	query, args := sb.Build()
	err := tr.db.Pg.QueryRow(query, args...).Scan(&total)
	if err != nil {
		return total, libs.DefaultInternalServerError(err)
	}

	return total, nil
}