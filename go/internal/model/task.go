package model

import (
	"time"

	"gopkg.in/guregu/null.v4"
)

type Task struct {
	ID int64 `json:"id,omitempty"`
	UserID int64 `json:"user_id,omitempty"`
	Title string `json:"title,omitempty"`
	Description *null.String `json:"description,omitempty"`
	Status string `json:"status,omitempty"`
	Order *null.Int `json:"order,omitempty"`
	CreatedAt *time.Time `json:"created_at,omitempty"`
	UpdatedAt *time.Time `json:"updated_at,omitempty"`
	Tags *[]Tag `json:"tags,omitempty"`
}