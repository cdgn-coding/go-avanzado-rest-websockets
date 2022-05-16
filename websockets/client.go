package websockets

import (
	"github.com/gorilla/websocket"
	"log"
)

type Client struct {
	hub      *Hub
	id       string
	socket   *websocket.Conn
	outbound chan []byte
}

func NewClient(hub *Hub, socket *websocket.Conn) *Client {
	return &Client{
		hub:      hub,
		socket:   socket,
		id:       socket.RemoteAddr().String(),
		outbound: make(chan []byte),
	}
}

func (c *Client) Write() {
	for {
		select {
		case message, ok := <-c.outbound:
			if !ok {
				err := c.socket.WriteMessage(websocket.CloseMessage, []byte{})
				if err != nil {
					log.Println("Error sending close connection message")
				}
				return
			}
			err := c.socket.WriteMessage(websocket.BinaryMessage, message)
			if err != nil {
				log.Println("Message lost")
			}
		}
	}
}
