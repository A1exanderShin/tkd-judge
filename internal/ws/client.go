package ws

import (
	"log"

	"github.com/gorilla/websocket"
)

type Client struct {
	conn *websocket.Conn
}

func NewClient(conn *websocket.Conn) *Client {
	return &Client{conn: conn}
}

func (c *Client) send(v any) {
	if err := c.conn.WriteJSON(v); err != nil {
		log.Printf("ws send error: %v", err)
	}
}

func (c *Client) sendRaw(data []byte) {
	if err := c.conn.WriteMessage(websocket.TextMessage, data); err != nil {
		log.Printf("ws send error: %v", err)
	}
}
