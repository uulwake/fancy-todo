package model

type User struct {
	ID int64 `json:"id,omitempty"`
	Name string `json:"name,omitempty"`
	Email string `json:"email,omitempty"`
	Password string `json:"password,omitempty"`
}