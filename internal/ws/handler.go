package ws

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

// Превращает HTTP-запрос в WebSocket
// CheckOrigin: return true - считается плохим тоном для прода, но при этом нормально в локальной сети
var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

type WSHandler struct {
	hub *Hub
}

func NewWSHandler(hub *Hub) *WSHandler {
	return &WSHandler{hub: hub}
}

func (h *WSHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		return
	}
	defer conn.Close()

	client := NewClient(conn)
	h.hub.register <- client
	defer func() {
		h.hub.unregister <- client
	}()

	for {
		_, data, err := conn.ReadMessage()
		if err != nil {
			log.Printf("ws read error: %v", err)
			return
		}

		var event Event
		if err := json.Unmarshal(data, &event); err != nil {
			log.Printf("invalid event json: %v", err)
			continue
		}

		h.hub.Publish(event)
	}
}
