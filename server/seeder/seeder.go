package seeder

import (
	"log"

	"github.com/RJD02/google-docs-clone/config"
)

var app config.AppConfig

func GetAppState(mainApp config.AppConfig) {
	app = mainApp
}

func insertUsers() {
	_, err := app.DBConn.Exec(`
    insert into users (username, email, password)
    values
    ('user1', 'user1@example.com', 'password1'),
    ('user2', 'user2@example.com', 'password2');
    `)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Inserted dummy data into users table")
}

func insertDocuments() {
	_, err := app.DBConn.Exec(`
    insert into documents (content, user_id, title)
    values
    ('Document 1 content', 1, 'Document 1'),
    ('Document 2 content', 2, 'Document 2');
`)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Inserted dummy data into the documents table.")
}

func insertUserDocumentRoles() {
	_, err := app.DBConn.Exec(`
    insert into user_document_role (user_id, document_id, permission)
    values
    (1, 1, 'OWNER'),
    (2, 2, 'VIEWER');
`)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Inserted dummy data into the user_document_role table.")
}

func RunSeeder() {
	insertUsers()
	insertDocuments()
	insertUserDocumentRoles()
}
