package ws

import (
	"log"

	"tkd-judge/internal/discipline"
)

type hubEvent struct {
	event  any
	client *Client
}

type Hub struct {
	discipline discipline.Discipline

	events chan hubEvent

	clients    map[*Client]struct{}
	mainJudge  *Client
	sideJudges map[int]*Client

	register   chan *Client
	unregister chan *Client
}

/* ================= CONSTRUCTOR ================= */

func NewHub(d discipline.Discipline) *Hub {
	return &Hub{
		discipline: d,

		events: make(chan hubEvent, 16),

		clients:    make(map[*Client]struct{}),
		sideJudges: make(map[int]*Client),

		register:   make(chan *Client),
		unregister: make(chan *Client),
	}
}

/* ================= CORE LOOP ================= */

func (h *Hub) Publish(e any, c *Client) {
	h.events <- hubEvent{event: e, client: c}
}

func (h *Hub) Run() {
	for {
		select {

		case ev := <-h.events:
			h.handleEvent(ev.event, ev.client)

		case c := <-h.register:
			h.handleRegister(c)

		case c := <-h.unregister:
			h.handleUnregister(c)
		}
	}
}

/* ================= EVENT ROUTER ================= */

func (h *Hub) handleEvent(event any, c *Client) {
	// роль-фильтр
	if !h.canSendEvent(c) {
		return
	}

	if err := h.discipline.HandleEvent(event); err != nil {
		log.Println("discipline error:", err)
		return
	}

	h.broadcastSnapshot()
}

/* ================= CLIENT MANAGEMENT ================= */

func (h *Hub) handleRegister(c *Client) {
	if c.role == RoleMainJudge {
		if h.mainJudge != nil {
			c.close()
			return
		}
		h.mainJudge = c
		log.Println("MAIN JUDGE CONNECTED")
	}

	if c.role == RoleSideJudge {
		h.sideJudges[c.judgeID] = c
		log.Printf("SIDE JUDGE %d CONNECTED", c.judgeID)
	}

	h.clients[c] = struct{}{}
	h.sendSnapshotTo(c)
}

func (h *Hub) handleUnregister(c *Client) {
	delete(h.clients, c)

	if c == h.mainJudge {
		h.mainJudge = nil
	}

	if c.role == RoleSideJudge {
		delete(h.sideJudges, c.judgeID)
	}
}

/* ================= ACCESS CONTROL ================= */

func (h *Hub) canSendEvent(c *Client) bool {
	if c.role == RoleMainJudge {
		return true
	}
	if c.role == RoleSideJudge {
		return true
	}
	return false
}

/* ================= BROADCAST ================= */

func (h *Hub) broadcastSnapshot() {
	snapshot := h.discipline.Snapshot()

	for c := range h.clients {
		c.send(snapshot)
	}
}

func (h *Hub) sendSnapshotTo(c *Client) {
	c.send(h.discipline.Snapshot())
}
