package ws

import (
	"judgement/internal/fight"
	"log"
)

type Hub struct {
	fight  *fight.Fight
	events chan Event
}

// NewHub создаёт Hub с новым боем
func NewHub() *Hub {
	return &Hub{
		fight:  fight.NewFight(),
		events: make(chan Event, 16),
	}
}

// Run запускает event-loop (блокирующий)
func (h *Hub) Run() {
	for event := range h.events {
		h.handleEvent(event)
	}
}

// Publish — единственная точка входа событий в Hub
func (h *Hub) Publish(event Event) {
	h.events <- event
}

// handleEvent — обработка события
func (h *Hub) handleEvent(event Event) {
	switch event.Type {
	case EventFightControl:
		h.handleFightControl(event.Data)
	default:
		log.Printf("unknown event type: %v", event.Type)
	}
}

func (h *Hub) handleFightControl(data any) {
	evt, ok := data.(FightControlEvent)
	if !ok {
		log.Printf("invalid fight control payload")
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
}
