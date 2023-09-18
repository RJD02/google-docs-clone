package handlers

import (
	"fmt"
	"log"
	"net/http"
	"sync"

	"github.com/RJD02/google-docs-clone/config"
	"github.com/RJD02/google-docs-clone/utils"
	"github.com/gofrs/uuid"
	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
)

var app config.AppConfig

type SocketDBInfo struct {
	count int
	rowId string
}

var documentConnectionsCount map[string]SocketDBInfo

func GetAppState(mainApp config.AppConfig) {
	app = mainApp
	documentConnectionsCount = make(map[string]SocketDBInfo)
}

func generateUniqueID() string {
	id, err := uuid.NewV4()
	if err != nil {
		log.Fatalln("Error generating UUID:", err)
		return ""
	}
	fmt.Println(id)
	return id.String()
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func HandleWebSocket(w http.ResponseWriter, r *http.Request) {
	wg := sync.WaitGroup{}
	vars := mux.Vars(r)
	documentId := vars["id"]

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer conn.Close()

	fmt.Println("Client connected")

	uniqueID := generateUniqueID()
	client := &utils.Client{
		ID:     uniqueID,
		Socket: conn,
	}
	app.WebSocketManager.RegisterClient(client)

	if _, exists := documentConnectionsCount[documentId]; !exists {
		documentConnectionsCount[documentId] = SocketDBInfo{
			count: 1,
		}
	} else {
		val, _ := documentConnectionsCount[documentId]
		val.count += 1
		documentConnectionsCount[documentId] = val
	}

	wg.Add(1)
	go conn.SetCloseHandler(func(code int, text string) error {
		defer wg.Done()
		defer conn.Close()
		if documentConnectionsCount[documentId].count == 1 {
			delete(documentConnectionsCount, documentId)
		} else {
			val, _ := documentConnectionsCount[documentId]
			val.count -= 1
			documentConnectionsCount[documentId] = val
		}
		log.Printf("Connection closed with code %d: %s", code, text)
		log.Println("Deleting the unique id associated with it")
		delete(app.WebSocketManager.Clients, uniqueID)
		return nil
	})

	// go Read(uniqueID, conn)

	wg.Add(1)
	// read incoming changes
	go writeIncomingChanges(&wg, conn, documentId, uniqueID)

	wg.Add(1)
	// subscribe to incoming changes for this document
	go subscribeToIncomingChanges(&wg, documentId, conn, uniqueID)
	wg.Wait()
}
