package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/RJD02/google-docs-clone/config"
	"github.com/RJD02/google-docs-clone/db"
	"github.com/RJD02/google-docs-clone/handlers"
	"github.com/RJD02/google-docs-clone/routes"
	"github.com/RJD02/google-docs-clone/seeder"
	"github.com/RJD02/google-docs-clone/types"
	"github.com/gorilla/mux"
	"gopkg.in/rethinkdb/rethinkdb-go.v6"
)

func main() {
	app := config.AppConfig{}
	app.SetEnvironment(config.DEVELOPMENT)
	app.SetShouldSeed(true)
	app.SetDBConnection(db.ConnectToDB())
	app.SetRethinkDBConnection(db.ConnectToRethinkDB())
	if app.Environment == config.DEVELOPMENT && app.ShouldSeed {
		go func() {
			db.GetAppState(app)
			seeder.GetAppState(app)
			seeder.RunSeeder()
			db.Run()
		}()
	}

	r := mux.NewRouter()

	r.HandleFunc("/ws", handlers.HandleWebSocket)
	go func() {
		fmt.Println("Goroutine started...")
		for {
			db.CreateRethinkDocumentTable()
			fmt.Println("Executed...")
			time.Sleep(5 * time.Second)
		}
	}()

	go func() {
		cursor, err := rethinkdb.Table("Document").Changes().Run(app.RethinkDBSess)
		if err != nil {
			log.Println("Error listening to changes:", err)
		}

		defer cursor.Close()

		var change types.RethinkChange

		for cursor.Next(&change) {
			fmt.Printf("Real-time change: %v\n", change)
		}
	}()

	serverAddr := "localhost:8000"
	log.Println("Websocket server started at ws://", serverAddr)

	http.Handle("/", r)
	http.Handle("/api", routes.ApiRouter)

	log.Println("Starting server in mode =", app.Environment)
	log.Println("Seeding =", app.ShouldSeed)

	err := http.ListenAndServe(serverAddr, nil)
	if err != nil {
		log.Println("Error starting WebSocket server:", err)
	}
}
