package websockets

import (
	"encoding/json"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
	"sync"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

type Hub struct {
	clients    []*Client
	register   chan *Client
	unregister chan *Client
	mutex      *sync.Mutex
}

func NewHub() *Hub {
	return &Hub{
		clients:    make([]*Client, 0),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		mutex:      &sync.Mutex{},
	}
}

func (hub *Hub) HandleWebSocket(w http.ResponseWriter, r *http.Request) {
	socket, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		http.Error(w, "Could not create websocket connection", http.StatusBadRequest)
		return
	}

	client := NewClient(hub, socket)
	hub.register <- client

	go client.Write()
}

func (hub *Hub) Run() {
	for {
		select {
		case client := <-hub.register:
			hub.onConnect(client)
		case client := <-hub.unregister:
			hub.onDisconnect(client)
		}
	}
}

func (hub *Hub) onConnect(client *Client) {
	log.Println("Client connected", client.id)
	hub.mutex.Lock()
	hub.clients = append(hub.clients, client)
	hub.mutex.Unlock()
}

func (hub *Hub) onDisconnect(client *Client) {
	log.Println("Client disconnected", client.id)
	err := client.socket.Close()
	if err != nil {
		log.Println("Error closing connection", client.id)
	}
	hub.mutex.Lock()
	var indexToRemove = -1
	for index, currentClient := range hub.clients {
		if client.id == currentClient.id {
			indexToRemove = index
		}
	}
	if indexToRemove != -1 {
		hub.clients = append(hub.clients[:indexToRemove], hub.clients[indexToRemove+1:]...)
	}
	hub.mutex.Unlock()
}

func (hub *Hub) Broadcast(message interface{}, ignore *Client) {
	data, _ := json.Marshal(message)
	for _, client := range hub.clients {
		if client != ignore {
			client.outbound <- data
		}
	}
}
