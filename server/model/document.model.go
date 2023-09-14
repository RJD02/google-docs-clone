package model

import (
	"time"
)

type Document struct {
	Id        int       `json:"id"`
	Content   string    `json:"content"`
	UserId    int       `json:"user_id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	ViewLink  string    `json:"view_link"`
	EditLink  string    `json:"edit_link"`
	IsPublic  bool      `json:"is_public"`
	Title     string    `json:"title"`
}
