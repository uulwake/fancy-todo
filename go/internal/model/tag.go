package model

import "time"

type Tag struct {
	ID int64 `json:"id,omitempty"`
	UserID int64 `json:"user_id,omitempty"`
	Name string `json:"name,omitempty"`
	CreatedAt time.Time `json:"created_at,omitempty"`
	UpdatedAt time.Time `json:"updated_at,omitempty"`
	Task Task `json:"task"`
}