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

func (tr *TagRepo) AddExistingTagToTask(ctx context.Context, userId int64, tagId int64, taskId int64) error {
	now := time.Now()
	sb := sqlbuilder.PostgreSQL.NewSelectBuilder()
	query, args := sb.
		Select("COUNT(id)").
		From("tasks").
		Where(sb.Equal("id", taskId), sb.Equal("user_id", userId)).
		Build()

	var totalTask int
	err := tr.db.Pg.QueryRow(query, args...).Scan(&totalTask)
	if err != nil {
		return libs.DefaultInternalServerError(err)
	}

	if totalTask != 1 {
		return libs.CustomError{
			HTTPCode: http.StatusBadRequest,
			Message: fmt.Sprintf("user ID: %d does not have task ID: %d", userId, taskId),
		}
	}

	sb = sqlbuilder.PostgreSQL.NewSelectBuilder()
	query, args = sb.
		Select("COUNT(id)").
		From("tags").
		Where(sb.Equal("id", taskId), sb.Equal("user_id", userId)).
		Build()

	var totalTag int
	err = tr.db.Pg.QueryRow(query, args...).Scan(&totalTag)
	if err != nil {
		return libs.DefaultInternalServerError(err)
	}

	if totalTask != 1 {
		return libs.CustomError{
			HTTPCode: http.StatusBadRequest,
			Message: fmt.Sprintf("user ID: %d does not have tag ID: %d", userId, tagId),
		}
	}
	
	ib := sqlbuilder.PostgreSQL.NewInsertBuilder()
	query, args = ib.InsertInto("tasks_tags").
		Cols("task_id", "tag_id", "created_at", "updated_at").
		Values(taskId, tagId, now, now).
		Build()

	_, err = tr.db.Pg.Exec(query, args...)
	if err != nil {
		return libs.DefaultInternalServerError(err)
	}

	return nil
}

func (tr *TagRepo) Search(ctx context.Context, userId int64, name string) ([]model.Tag, error) {
	tags := []model.Tag{}

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
						"name": map[string]any {
							"value": fmt.Sprintf("*%s*", name),
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
		return tags, libs.DefaultInternalServerError(err)
	}

	response, err := tr.db.Es.Search(
		tr.db.Es.Search.WithIndex(constant.ES_INDEX_TAGS),
		tr.db.Es.Search.WithBody(&bodyJson),
	)
	if err != nil {
		return tags, libs.DefaultInternalServerError(err)
	}
	if response.StatusCode != http.StatusOK {
		return tags, libs.CustomError{
			HTTPCode: http.StatusBadRequest,
			Message: "error when searching tags in ElasticSearch",
		}
	}


	result := make(map[string]any)
	json.NewDecoder(response.Body).Decode(&result)

	hits, ok := result["hits"].(map[string]any)["hits"]
	if !ok {
		return tags, libs.CustomError{
			HTTPCode: http.StatusInternalServerError,
			Message: "invalid es response",
		}
	}

	for _, hit := range hits.([]any) {
		tag := model.Tag{}
		source := hit.(map[string]any)["_source"].(map[string]any)

		tag.ID = int64(source["id"].(float64))
		tag.UserID = int64(source["user_id"].(float64))
		tag.Name = source["name"].(string)

		tags = append(tags, tag)
	}

	return tags, nil
}

func (tr *TagRepo) DeleteById(ctx context.Context, userId int64, tagId int64) error {
	db := sqlbuilder.PostgreSQL.NewDeleteBuilder()
	db.DeleteFrom("tags")
	db.Where(db.Equal("id", tagId), db.Equal("user_id", userId))

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
							"id": tagId,
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

	response, err := tr.db.Es.DeleteByQuery([]string{constant.ES_INDEX_TAGS}, &bodyJson)
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