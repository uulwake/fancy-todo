package model

import "time"

type Task struct {
	ID int64 `json:"id,omitempty"`
	UserID int64 `json:"user_id,omitempty"`
	Title string `json:"title,omitempty"`
	Description string `json:"description,omitempty"`
	Status string `json:"status,omitempty"`
	Order int `json:"order,omitempty"`
	CreatedAt time.Time `json:"created_at,omitempty"`
	UpdatedAt time.Time `json:"updated_at,omitempty"`
}