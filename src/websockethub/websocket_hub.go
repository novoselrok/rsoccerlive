package websockethub

import (
	"github.com/gorilla/websocket"
)

type Client struct {
	*websocket.Conn
}

type Hub struct {
	Clients map[*Client]bool
}

func (hub *Hub) RegisterClient(client *Client) {
	hub.Clients[client] = true
}

func (hub *Hub) BroadcastJSON(message interface{}) {
	for client := range hub.Clients {
		err := client.WriteJSON(message)
		if err != nil {
			client.Close()
			delete(hub.Clients, client)
		}
	}
}

func NewWebsocketHub() *Hub {
	return &Hub{
		Clients: make(map[*Client]bool),
	}
}
