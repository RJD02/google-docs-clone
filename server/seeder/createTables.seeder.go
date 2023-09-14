package seeder

import (
	"fmt"

	"github.com/RJD02/google-docs-clone/config"
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
    )
        `
	_, err := app.DBConn.Exec(createUserTableQuery)
	if err != nil {
		panic(err.Error())
	}
	fmt.Println("Table `users` created successfully")
}
