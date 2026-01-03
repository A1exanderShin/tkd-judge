package ws

import (
	"encoding/json"
	"log"
	"time"
	"tkd-judge/internal/events"
	"tkd-judge/internal/judges"

	"tkd-judge/internal/fight"
)

type Hub struct {
	fight *fight.Fight

	scoreboard *fight.Scoreboard
	judges     map[int]*judges.Judge
	eventLog   []events.ScoreEvent

	events chan Event

	clients    map[*Client]struct{}
	register   chan *Client
	unregister chan *Client
}

func NewHub() *Hub {
	j := make(map[int]*judges.Judge)
	for i := 1; i <= 4; i++ {
		j[i] = judges.NewJudge(i, 300*time.Millisecond)
	}

	return &Hub{
		fight:      fight.NewFight(),
		scoreboard: fight.NewScoreboard(),
		judges:     j,
		eventLog:   make([]events.ScoreEvent, 0),
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
	case EventScore:
		h.handleScore(event.Data)
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

func (h *Hub) handleScore(data json.RawMessage) {
	// 1. бой должен идти
	if h.fight.State() != fight.StateRunning {
		log.Println("score ignored: fight not running")
		return
	}

	var payload ScorePayload
	if err := json.Unmarshal(data, &payload); err != nil {
		log.Printf("invalid score payload: %v", err)
		return
	}

	judge, ok := h.judges[payload.JudgeID]
	if !ok {
		log.Printf("unknown judge: %d", payload.JudgeID)
		return
	}

	now := time.Now()
	if err := judge.CanScore(now); err != nil {
		log.Printf("judge %d click ignored: %v", payload.JudgeID, err)
		return
	}

	event := events.ScoreEvent{
		JudgeID: payload.JudgeID,
		Fighter: events.Fighter(payload.Fighter),
		Points:  payload.Points,
		Time:    now,
	}

	h.eventLog = append(h.eventLog, event)
	h.scoreboard.Apply(event)

	red, blue := h.scoreboard.Score()
	log.Printf(
		"SCORE: judge=%d fighter=%s +%d | TOTAL red=%d blue=%d",
		event.JudgeID, event.Fighter, event.Points, red, blue,
	)

	h.broadcastScore()
}

func (h *Hub) broadcastScore() {
	red, blue := h.scoreboard.Score()

	msg := map[string]any{
		"type": "score_update",
		"red":  red,
		"blue": blue,
	}

	for client := range h.clients {
		client.send(msg)
	}
}
