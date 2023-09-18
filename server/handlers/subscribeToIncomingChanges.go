package handlers

import (
	"log"
	"sync"

	"github.com/gorilla/websocket"
	"gopkg.in/rethinkdb/rethinkdb-go.v6"
)

func subscribeToIncomingChanges(wg *sync.WaitGroup, documentId string, conn *websocket.Conn, uniqueId string) {
	defer wg.Done()
	filter := map[string]interface{}{
		"DocumentId": documentId,
		"Id":         uniqueId,
	}
	cursor, err := rethinkdb.Table("Document").Filter(filter).Changes().Run(app.RethinkDBSess)
	if err != nil {
		log.Println("Error listening to changes", err)
		return
	}
	defer cursor.Close()
	var change map[string]map[string]interface{}
	for cursor.Next(&change) {
		conn.WriteJSON(change)
	}
}
