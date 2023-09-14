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

var id = 1
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

	conn.SetCloseHandler(func(code int, text string) error {
		log.Printf("Connection closed with code %d: %s", code, text)
		return nil
	})

	fmt.Println("Client connected")

	uniqueID := generateUniqueID()
	clients[uniqueID] = &WebSocketClient{
		ID:   uniqueID,
		Conn: conn,
	}

	for {
		fmt.Println(len(clients))
		messageType, p, err := conn.ReadMessage()
		if err != nil {
			fmt.Println(err)
			return
		}

		log.Printf("recv: %d", messageType)
		log.Printf("msg: %s\n", p)
		log.Println("Generated uuid: ", clients[uniqueID].ID)

		returnMsg := "Hello"
		if messageType == websocket.TextMessage {
			returnMsg += clients[uniqueID].ID
		}

		if err := conn.WriteMessage(messageType, []byte(returnMsg)); err != nil {
			fmt.Println(err)
			return
		}
	}
}
