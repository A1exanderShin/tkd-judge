package ws

import (
	"encoding/json"
	"log"

	"tkd-judge/internal/fight"
)

type Hub struct {
	fight *fight.Fight

	events chan Event

	clients    map[*Client]struct{}
	register   chan *Client
	unregister chan *Client
}

func NewHub() *Hub {
	return &Hub{
		fight:      fight.NewFight(),
		events:     make(chan Event, 16),
		clients:    make(map[*Client]struct{}),
		register:   make(chan *Client),
		unregister: make(chan *Client),
	}
}

func (h *Hub) Run() {
	for {
		select {
		case event := <-h.events:
			h.handleEvent(event)

		case client := <-h.register:
			h.clients[client] = struct{}{}

			// 🔑 отправляем текущее состояние сразу при подключении
			client.send(map[string]string{
				"type":  "state",
				"state": h.fight.State().String(),
			})

		case client := <-h.unregister:
			delete(h.clients, client)
		}
	}
}

func (h *Hub) Publish(event Event) {
	h.events <- event
}

func (h *Hub) handleEvent(event Event) {
	switch event.Type {
	case EventFightControl:
		h.handleFightControl(event.Data)
	default:
		log.Printf("unknown event type: %v", event.Type)
	}
}

func (h *Hub) handleFightControl(data json.RawMessage) {
	var evt FightControlEvent

	if err := json.Unmarshal(data, &evt); err != nil {
		log.Printf("invalid fight control payload: %v", err)
		return
	}

	var err error

	switch evt.Action {
	case ActionStart:
		err = h.fight.Start()
	case ActionPause:
		err = h.fight.Pause()
	case ActionResume:
		err = h.fight.Resume()
	case ActionStop:
		err = h.fight.Stop()
	default:
		log.Printf("unknown fight action: %v", evt.Action)
		return
	}

	if err != nil {
		log.Printf("fight action error: %v", err)
		return
	}

	log.Printf("fight state changed to %s", h.fight.State())
	h.broadcastState()
}

func (h *Hub) broadcastState() {
	msg := map[string]string{
		"type":  "state",
		"state": h.fight.State().String(),
	}

	for client := range h.clients {
		client.send(msg)
	}
}
