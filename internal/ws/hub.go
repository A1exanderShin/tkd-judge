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
	// 🔥 подписка на realtime события дисциплины
	go h.listenRealtime()

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

/* ================= REALTIME ================= */

func (h *Hub) listenRealtime() {
	for ev := range h.discipline.Realtime() {
		for c := range h.clients {
			c.send(map[string]any{
				"type": ev.Type,
				"data": ev.Data,
			})
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

	// snapshot шлём ТОЛЬКО после доменных изменений
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

		// отправляем judge_id один раз
		c.send(map[string]any{
			"type":    "judge_id",
			"judgeID": c.judgeID,
		})
	}

	h.clients[c] = struct{}{}

	// при подключении сразу шлём snapshot
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
	switch c.role {
	case RoleMainJudge, RoleSideJudge:
		return true
	default:
		return false
	}
}

/* ================= SNAPSHOT ================= */

func (h *Hub) broadcastSnapshot() {
	payload := map[string]any{
		"type": "snapshot",
		"data": h.discipline.Snapshot(),
	}

	for c := range h.clients {
		c.send(payload)
	}
}

func (h *Hub) sendSnapshotTo(c *Client) {
	c.send(map[string]any{
		"type": "snapshot",
		"data": h.discipline.Snapshot(),
	})
}
