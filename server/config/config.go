package config

import "database/sql"

const (
	DEVELOPMENT string = "DEVELOPMENT"
	PRODUCTION         = "PRODUCTION"
)

type AppConfig struct {
	Environment string
	DBConn      *sql.DB
	ShouldSeed  bool
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

func (app *AppConfig) SetShouldSeed(val bool) {
	app.ShouldSeed = val
}
