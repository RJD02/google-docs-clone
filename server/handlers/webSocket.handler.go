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
	"gopkg.in/rethinkdb/rethinkdb-go.v6"
)

var app config.AppConfig

func GetAppState(mainApp config.AppConfig) {
	app = mainApp
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

func writeIncomingChanges(wg *sync.WaitGroup, conn *websocket.Conn, documentId string, uniqueID string) {

	defer wg.Done()
	for {
		fmt.Println(len(app.WebSocketManager.Clients))
		messageType, p, err := conn.ReadMessage()
		if err != nil {
			fmt.Println(err)
			return
		}

		log.Printf("recv: %d", messageType)
		log.Printf("msg: %s\n", p)
		log.Println("Generated uuid: ", app.WebSocketManager.Clients[uniqueID].ID)

		documentData := make(map[string]string)

		documentData["DocumentId"] = documentId
		documentData["Id"] = uniqueID
		documentData["Content"] = string(p)

		rethinkdb.Table("Document").Insert(documentData).RunWrite(app.RethinkDBSess)

		returnMsg := string(p)

		if err := conn.WriteMessage(messageType, []byte(returnMsg)); err != nil {
			fmt.Println(err)
			return
		}
	}
}

func subscribeToIncomingChanges(wg *sync.WaitGroup, documentId string, conn *websocket.Conn, uniqueId string) {
	defer wg.Done()
	filter := map[string]interface{}{
		"DocumentId": documentId,
	}
	log.Printf("%v", filter)
	cursor, err := rethinkdb.Table("Document").Filter(filter).Changes().Run(app.RethinkDBSess)
	if err != nil {
		log.Println("Error listening to changes", err)
		return
	}
	defer cursor.Close()
	var change map[string]map[string]interface{}
	for cursor.Next(&change) {
		log.Println("Running because of change", uniqueId)
		id := change["new_val"]["Id"]
		if str, ok := id.(string); ok {
			// conversion can be done
			_, exists := app.WebSocketManager.Clients[str]
			if !exists {
				log.Println("No such id exists in current pool of websockets")
				return
			}
			changedContent := change["new_val"]["Content"]
			log.Println(changedContent, uniqueId)
			if byt, ok := changedContent.([]byte); ok {
				if err := conn.WriteMessage(websocket.TextMessage, byt); err != nil {
					log.Fatal(err)
					return
				} else {
					log.Println("Message sent")
				}
			}

		}
	}
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

	wg.Add(1)
	go conn.SetCloseHandler(func(code int, text string) error {
		defer wg.Done()
		defer conn.Close()
		log.Printf("Connection closed with code %d: %s", code, text)
		log.Println("Deleting the unique id associated with it")
		delete(app.WebSocketManager.Clients, uniqueID)
		_, exists := app.WebSocketManager.Clients[uniqueID]
		if !exists {
			log.Println("Successfully deleted the unique key")
		}
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

func ListenToRethinkDB() {
	cursor, err := rethinkdb.Table("Document").Changes().Run(app.RethinkDBSess)

	if err != nil {
		log.Println("Error Listening to changes", err)
		return

	}
	defer cursor.Close()
	var change map[string]map[string]interface{}

	for cursor.Next(&change) {
		id := change["new_val"]["Id"]
		if str, ok := id.(string); ok {
			// conversion can be done
			_, exists := app.WebSocketManager.Clients[str]
			if !exists {
				log.Println("No such id exists in current pool of websockets")
			}
			changedContent := change["new_val"]["Content"]
			if byt, ok := changedContent.([]byte); ok {
				if err = app.WebSocketManager.Clients[str].Socket.WriteMessage(websocket.TextMessage, []byte(byt)); err != nil {
					log.Fatal(err)
					return
				}
			}
		}
	}
}
