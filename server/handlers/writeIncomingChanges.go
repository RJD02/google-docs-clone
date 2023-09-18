package handlers

import (
	"fmt"
	"log"
	"sync"

	"github.com/gorilla/websocket"
	"gopkg.in/rethinkdb/rethinkdb-go.v6"
)

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

		if _, exists := documentConnectionsCount[documentId]; exists {
			rethinkdb.Table("Document").Get(documentConnectionsCount[documentId].rowId).Update(documentData).Run(app.RethinkDBSess)
		} else {
			response, err := rethinkdb.Table("Document").Insert(documentData).RunWrite(app.RethinkDBSess)
			if err != nil {
				log.Fatal(err)
			}
			insertedID := response.GeneratedKeys[0]
			documentConnectionsCount[documentId] = SocketDBInfo{
				count: 1,
				rowId: insertedID,
			}
		}

		returnMsg := string(p)

		conn.WriteJSON(returnMsg)
	}
}
