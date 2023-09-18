package handlers

import (
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/gorilla/websocket"
	"gopkg.in/rethinkdb/rethinkdb-go.v6"
)

func writeIncomingChanges(wg *sync.WaitGroup, conn *websocket.Conn, documentId string, uniqueID string) {

	defer wg.Done()
	for {
		messageType, p, err := conn.ReadMessage()
		if err != nil {
			fmt.Println(err, messageType)
			return
		}

		documentData := make(map[string]string)

		documentData["DocumentId"] = documentId
		documentData["Id"] = uniqueID
		documentData["Content"] = string(p)
		documentData["timestamp"] = time.Now().String()

		if _, exists := documentConnectionsCount[documentId]; exists {
			log.Println("Document exists")
			rowId := documentConnectionsCount[documentId].rowId
			rethinkdb.Table("Document").Get(rowId).Update(documentData).Run(app.RethinkDBSess)
		} else {
			response, err := rethinkdb.Table("Document").Insert(documentData).RunWrite(app.RethinkDBSess)
			if err != nil {
				log.Fatal(err)
			}
			insertedID := response.GeneratedKeys[0]
			log.Println("Inserted with id", insertedID)
			documentConnectionsCount[documentId] = SocketDBInfo{
				count: 1,
				rowId: insertedID,
			}
		}

		returnMsg := string(p)

		conn.WriteJSON(returnMsg)
	}
}
