package main

import (
	"log"
	"net/http"

	"github.com/RJD02/google-docs-clone/config"
	"github.com/RJD02/google-docs-clone/db"
	"github.com/RJD02/google-docs-clone/handlers"
	"github.com/RJD02/google-docs-clone/routes"
	"github.com/RJD02/google-docs-clone/seeder"
	"github.com/gorilla/mux"
)

func main() {
	app := config.AppConfig{}
	app.SetEnvironment(config.DEVELOPMENT)
	app.SetShouldSeed(false)
	app.SetDBConnection(db.ConnectToDB())
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
