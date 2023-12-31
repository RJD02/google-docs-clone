package config

import (
	"database/sql"

	"github.com/RJD02/google-docs-clone/utils"
	rethinkdb "gopkg.in/rethinkdb/rethinkdb-go.v6"
)

const (
	DEVELOPMENT string = "DEVELOPMENT"
	PRODUCTION         = "PRODUCTION"
)

type AppConfig struct {
	Environment        string
	DBConn             *sql.DB
	ShouldSeed         bool
	RethinkDBSess      *rethinkdb.Session
	DBTableName        string
	RethinkDBTableName string
	WebSocketManager   *utils.WebSocketManager
}

func (app *AppConfig) SetWebSocketManager(manager *utils.WebSocketManager) {
	app.WebSocketManager = manager

}

func (app *AppConfig) SetEnvironment(environment string) {
	if environment == DEVELOPMENT || environment == PRODUCTION {
		app.Environment = environment
	} else {
		app.Environment = DEVELOPMENT
	}
}

func (app *AppConfig) SetDBConnection(db *sql.DB) {
	app.DBConn = db
}

func (app *AppConfig) SetRethinkDBConnection(sess *rethinkdb.Session) {
	app.RethinkDBSess = sess
}

func (app *AppConfig) SetShouldSeed(val bool) {
	app.ShouldSeed = val
}
