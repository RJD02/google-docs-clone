package model

const (
	ADMIN  string = "ADMIN"
	VIEWER        = "VIEWER"
	EDITOR        = "EDITOR"
	OWNER         = "OWNER"
)

type UserDocumentRole struct {
	UserId     int    `json:"user_id"`
	DocumentId int    `json:"document_id"`
	Permission string `json:"permission"`
	Id         int    `json:"id"`
}
