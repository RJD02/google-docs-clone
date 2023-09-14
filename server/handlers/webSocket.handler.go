package handlers

import (
	"fmt"
	"github.com/gofrs/uuid"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
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

func HandleWebSocket(w http.ResponseWriter, r *http.Request) {
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
