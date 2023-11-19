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
	"net/http"
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
	
	ib := sqlbuilder.PostgreSQL.NewInsertBuilder()
	query, args := ib.InsertInto("tasks").
		Cols("user_id", "title", "description", "status", "created_at", "updated_at").
		Values(data.UserId, data.Title, data.Description,  constant.TASK_STATUS_ON_GOING,now, now).
		SQL("RETURNING id").
		Build()

	err = tx.QueryRow(query, args...).Scan(&taskId)
	if err != nil {
		return 0, libs.DefaultInternalServerError(err)
	}

	if len(data.TagIDs) > 0 {
		ib := sqlbuilder.PostgreSQL.NewInsertBuilder()
		ib.InsertInto("tasks_tags").Cols("task_id", "tag_id", "created_at", "updated_at")

		for _, tagId := range data.TagIDs {
			ib.Values(taskId, tagId, now, now)
		}

		query, args := ib.Build()
		_, err = tx.Exec(query, args...)
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
	
	var taskJson bytes.Buffer
	err = json.NewEncoder(&taskJson).Encode(&task)
	if err != nil {
		return 0, libs.DefaultInternalServerError(err)
	}

	response, err := tr.db.Es.Index(constant.ES_INDEX_TASKS, &taskJson)
	if err != nil {
		return 0, libs.DefaultInternalServerError(err)
	}
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

func (tr *TaskRepo) Search(ctx context.Context, userId int64, title string) ([]model.Task, error) {
	tasks := []model.Task{}

	body := map[string]any{
		"query": map[string]any {
			"bool": map[string]any{
				"must": map[string]any{
					"match": map[string]int64 {
						"user_id": userId,
					},
				},
				"filter": map[string]any {
					"wildcard": map[string]any {
						"title": map[string]any {
							"value": fmt.Sprintf("*%s*", title),
							"case_insensitive": true,
						},
					},
				},
			},
		},
	}

	var bodyJson bytes.Buffer
	err := json.NewEncoder(&bodyJson).Encode(body)
	if err != nil {
		return tasks, libs.DefaultInternalServerError(err)
	}

	response, err := tr.db.Es.Search(
		tr.db.Es.Search.WithIndex(constant.ES_INDEX_TASKS),
		tr.db.Es.Search.WithBody(&bodyJson),
	)
	if err != nil {
		return tasks, libs.DefaultInternalServerError(err)
	}
	if response.StatusCode != http.StatusOK {
		return tasks, libs.CustomError{
			HTTPCode: http.StatusBadRequest,
			Message: "error when searching tags in ElasticSearch",
		}
	}

	result := make(map[string]any)
	json.NewDecoder(response.Body).Decode(&result)

	hits, ok := result["hits"].(map[string]any)["hits"]
	if !ok {
		return tasks, libs.CustomError{
			HTTPCode: http.StatusInternalServerError,
			Message: "invalid es response",
		}
	}

	for _, hit := range hits.([]any) {
		task := model.Task{}
		source := hit.(map[string]any)["_source"].(map[string]any)

		task.ID = int64(source["id"].(float64))
		task.UserID = int64(source["user_id"].(float64))
		task.Title = source["title"].(string)
		task.Status = source["status"].(string)

		tasks = append(tasks, task)
	}

	return tasks, nil
}

func (tr *TaskRepo) UpdateById(ctx context.Context, userId int64, taskId int64, task TaskUpdateByIdInput) error {
	ub := sqlbuilder.PostgreSQL.NewUpdateBuilder().Update("tasks")

	now := time.Now()
	ub.Set(ub.Assign("updated_at", now))
	esSources := []string{fmt.Sprintf("ctx._source.updated_at = '%s'", now.Format(time.RFC3339))}

	if task.Title != "" {
		ub.SetMore(ub.Assign("title", task.Title))
		esSources = append(esSources, fmt.Sprintf("ctx._source.title = '%s'", task.Title))
	}

	if task.Description != "" {
		ub.SetMore(ub.Assign("description", task.Description))
		esSources = append(esSources, fmt.Sprintf("ctx._source.description = '%s'", task.Description))
	}

	if task.Status != "" {
		ub.SetMore(ub.Assign("status", task.Status))
		esSources = append(esSources, fmt.Sprintf("ctx._source.status = '%s'", task.Status))
	}

	if task.Order != 0 {
		ub.SetMore(ub.Assign("\"order\"", task.Order))
		esSources = append(esSources, fmt.Sprintf("ctx._source.order = %d", task.Order))
	}

	query, args := ub.Where(ub.Equal("id", taskId)).Where(ub.Equal("user_id", userId)).Build()

	tx, err := tr.db.Pg.Begin()
	if err != nil {
		return libs.DefaultInternalServerError(err)
	}
	defer tx.Rollback()

	_, err = tx.Exec(query, args...)
	if err != nil {
		return libs.DefaultInternalServerError(err)
	}

	body := map[string]any {
		"query": map[string]any {
			"bool": map[string]any {
				"must": []map[string]any{
					{
						"match": map[string]any{
							"id": taskId,
						},
					},
					{
						"match": map[string]any{
							"user_id": userId,
						},
					},
				},
			},
		},
		"script": map[string]any{
			"lang": "painless",
			"source": strings.Join(esSources, ";"),
		},
	}

	var bodyJson bytes.Buffer
	err = json.NewEncoder(&bodyJson).Encode(body)
	if err != nil {
		return libs.DefaultInternalServerError(err)
	}

	response, err := tr.db.Es.UpdateByQuery(
		[]string{constant.ES_INDEX_TASKS}, 
		tr.db.Es.UpdateByQuery.WithBody(&bodyJson),
	)
	if err != nil {
		return libs.DefaultInternalServerError(err)
	}

	if response.StatusCode != http.StatusOK {
		return libs.CustomError{
			HTTPCode: http.StatusBadRequest,
			Message: "error updating data to ElasticSearch",
		}
	}

	if err := tx.Commit(); err != nil {
		return libs.DefaultInternalServerError(err)
	}

	return nil
}

func (tr *TaskRepo) DeleteById(ctx context.Context, userId int64, taskId int64) error {
	db := sqlbuilder.PostgreSQL.NewDeleteBuilder()
	db.DeleteFrom("tasks")
	db.Where(db.Equal("id", taskId), db.Equal("user_id", userId))

	query, args := db.Build()
	tx, err := tr.db.Pg.Begin()
	if err != nil {
		return libs.DefaultInternalServerError(err)
	}
	defer tx.Rollback()

	_, err = tx.Exec(query, args...)
	if err != nil {
		return libs.DefaultInternalServerError(err)
	}

	body := map[string]any {
		"query": map[string]any {
			"bool": map[string]any {
				"must": []map[string]any {
					{
						"match": map[string]any {
							"id": taskId,
						},
					},
					{
						"match": map[string]any {
							"user_id": userId,
						},
					},
				},
			},
		},
	}

	var bodyJson bytes.Buffer
	err = json.NewEncoder(&bodyJson).Encode(body)
	if err != nil {
		return libs.DefaultInternalServerError(err)
	}

	response, err := tr.db.Es.DeleteByQuery([]string{constant.ES_INDEX_TASKS}, &bodyJson)
	if err != nil {
		return libs.DefaultInternalServerError(err)
	}

	if response.StatusCode != http.StatusOK {
		return libs.CustomError{
			HTTPCode: http.StatusBadRequest,
			Message: "error deleting data to ElasticSearch",
		}
	}

	if err = tx.Commit(); err != nil {
		return libs.DefaultInternalServerError(err)
	}

	return nil
}