package utils

import "github.com/gorilla/websocket"

type Client struct {
	ID         string
	Socket     *websocket.Conn
	DocumentId string
}

type WebSocketManager struct {
	Clients map[string]*Client
}

func NewWebSocketManager() *WebSocketManager {
	return &WebSocketManager{
		Clients: make(map[string]*Client),
	}
}

func (manager *WebSocketManager) RegisterClient(client *Client) {
	manager.Clients[client.ID] = client
}
