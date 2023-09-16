package db

import (
	"database/sql"
	"log"

	_ "github.com/go-sql-driver/mysql"
	rethinkdb "gopkg.in/rethinkdb/rethinkdb-go.v6"
)

const DB_NAME = "google_docs_clone"

func ConnectToDB() *sql.DB {
	log.Println("Connecting to mysql...")
	db, err := sql.Open("mysql", "Dulange:[X2JnrKiV[RbrRRU@tcp(127.0.0.1:3306)/"+DB_NAME)
	if err != nil {
		panic(err.Error())
	}
	log.Println("Connected to mysql")

	return db
}

func ConnectToRethinkDB() *rethinkdb.Session {
	session, err := rethinkdb.Connect(rethinkdb.ConnectOpts{
		Address:  "localhost",
		Database: "google_docs_clone",
	})

	if err != nil {
		panic(err)
	}
	return session
}
