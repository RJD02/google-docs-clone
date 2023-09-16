package db

import (
	"fmt"
	"log"
	"time"

	"github.com/RJD02/google-docs-clone/config"
	"github.com/RJD02/google-docs-clone/model"
	rethinkdb "gopkg.in/rethinkdb/rethinkdb-go.v6"
)

var app config.AppConfig

func GetAppState(mainApp config.AppConfig) {
	app = mainApp
}

func CreateUserTable() {
	createUserTableQuery := `
    create table if not exists users (
        id int auto_increment primary key,
        username varchar(255) not null,
        email varchar(200) not null,
        password varchar(20) not null,
        is_verfied boolean default 0,
        reset_password_token varchar(200),
        verification_token varchar(200)
    );
        `
	_, err := app.DBConn.Exec(createUserTableQuery)
	if err != nil {
		panic(err.Error())
	}
	log.Println("Table `users` created successfully")
}

func CreateDocumentTable() {
	CreateDocumentTableQuery := `
    create table if not exists documents (
        id int auto_increment primary key,
        content text,
        user_id int not null,
        created_at timestamp default current_timestamp,
        updated_at datetime default current_timestamp on update current_timestamp,
        view_link varchar(200),
        edit_link varchar(200),
        is_public boolean default 0,
        title varchar(100),
        foreign key (user_id) references users(id) on delete cascade
    );
    `
	_, err := app.DBConn.Exec(CreateDocumentTableQuery)
	if err != nil {
		panic(err.Error())
	}
	log.Println("Table `documents` created successfully")
}

func CreateUserDocumentRoleTable() {
	CreateUserDocumentRoleTableQuery := `
    create table if not exists user_document_role (
        id int auto_increment primary key,
        user_id int not null,
        document_id int not null,
        permission enum('ADMIN', 'VIEWER', 'EDITOR', 'OWNER'),
        foreign key (user_id) references users(id) on delete cascade,
        foreign key (document_id) references documents(id)
        on delete cascade
    );
    `
	_, err := app.DBConn.Exec(CreateUserDocumentRoleTableQuery)
	if err != nil {
		panic(err.Error())
	}
	log.Println("Table `user-document-role` created successfully")
}

func CreateRethinkDocumentTable() {
	document := model.Document{
		Id:        1,
		Content:   "Hello World, this is the first document",
		UserId:    1,
		Title:     "Something Interesting",
		IsPublic:  false,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	res, err := rethinkdb.Table("Document").Insert(document).RunWrite(app.RethinkDBSess)
	if err != nil {
		panic(err)
	}
	fmt.Println(res.GeneratedKeys[0])
}

func Run() {
	CreateUserTable()
	CreateDocumentTable()
	CreateUserDocumentRoleTable()
	CreateRethinkDocumentTable()
}
