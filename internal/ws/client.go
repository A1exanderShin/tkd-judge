package ws

import (
	"log"

	"github.com/gorilla/websocket"
)

// Client — тонкая обертка над websocket-соединением.
// Никакой логики, только отправка и закрытие.
type Client struct {
	conn    *websocket.Conn
	role    Role
	judgeID int // только для боковых судей (1–4)
}

func NewClient(conn *websocket.Conn, role Role) *Client {
	return &Client{
		conn: conn,
		role: role,
	}
}

func (c *Client) send(v any) {
	if err := c.conn.WriteJSON(v); err != nil {
		log.Printf("ws send error: %v", err)
	}
}

func (c *Client) close() {
	_ = c.conn.Close()
}
