package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/RJD02/google-docs-clone/config"
	"github.com/RJD02/google-docs-clone/db"
	"github.com/RJD02/google-docs-clone/seeder"
	"github.com/gofrs/uuid"
	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
)

type WebSocketClient struct {
	ID   string
	Conn *websocket.Conn
}

var clients = make(map[string]*WebSocketClient)

func generateUniqueID() string {
	id, err := uuid.NewV4()
	if err != nil {
		log.Fatalln("Error generating UUID:", err)
		return ""
	}
	fmt.Println(id)
	return ""
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func handleWebSocket(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer conn.Close()

	fmt.Println("Client connected")

	uniqueID := generateUniqueID()
	client := &WebSocketClient{
		ID:   uniqueID,
		Conn: conn,
	}
	clients[uniqueID] = client

	for {
		fmt.Println(len(clients))
		messageType, p, err := conn.ReadMessage()
		if err != nil {
			fmt.Println(err)
			return
		}

		log.Printf("recv: %d", messageType)
		log.Printf("msg: %s", p)

		if err := conn.WriteMessage(messageType, p); err != nil {
			fmt.Println(err)
			return
		}
	}
}

func main() {
	app := config.AppConfig{}
	app.SetEnvironment(config.DEVELOPMENT)
	fmt.Println(app.Environment)
	app.SetDBConnection(db.ConnectToDB())
	if app.Environment == config.DEVELOPMENT {
		seeder.GetAppState(app)
		seeder.CreateUserTable()
	}

	r := mux.NewRouter()
	r.HandleFunc("/ws", handleWebSocket)

	serverAddr := "localhost:8000"
	fmt.Println("Websocket server started at ws://", serverAddr)

	http.Handle("/", r)

	err := http.ListenAndServe(serverAddr, nil)
	if err != nil {
		fmt.Println("Error starting WebSocket server:", err)
	}
}
