package ws

import (
	"encoding/json"
	"log"
	"time"

	"tkd-judge/internal/config"
	"tkd-judge/internal/events"
	"tkd-judge/internal/fight"
	"tkd-judge/internal/judges"
)

type Hub struct {
	cfg config.Config

	fight *fight.Fight // FSM боя
	timer *fight.Timer

	scoreboard *fight.Scoreboard
	warnings   *fight.WarningCounter

	judges   map[int]*judges.Judge
	eventLog []any

	events chan Event // внешние события, один вход в систему, один обработчик (actor-model)

	// Hub из пакета Gorilla WS
	clients    map[*Client]struct{}
	register   chan *Client
	unregister chan *Client
}

func NewHub() *Hub {
	cfg := config.Default()

	// инициализация судей
	j := make(map[int]*judges.Judge)
	for i := 1; i <= cfg.JudgesCount; i++ {
		j[i] = judges.NewJudge(i, cfg.AntiClick)
	}

	// таймер
	timer := fight.NewTimer(cfg.RoundDuration)

	// собираем hub - вся система и все компоненты в одно целое
	h := &Hub{
		cfg:        cfg,
		fight:      fight.NewFight(),
		timer:      timer,
		scoreboard: fight.NewScoreboard(),
		warnings:   fight.NewWarningCounter(cfg.WarningsForPenalty),
		judges:     j,
		eventLog:   make([]any, 0),
		events:     make(chan Event, 16),
		clients:    make(map[*Client]struct{}),
		register:   make(chan *Client),
		unregister: make(chan *Client),
	}

	// тикание таймера, отправка в консоль
	timer.OnTick(func(rem time.Duration) {
		h.broadcastTimer(int(rem.Seconds()))
	})

	// таймер закончился, бой остановился, отправка состояния
	timer.OnFinished(func() {
		_ = h.fight.Stop()
		h.broadcastState()
	})

	return h
}

// event loop
func (h *Hub) Run() {
	for {
		select {
		// входящие события, все изменения состояния здесь
		case event := <-h.events:
			h.handleEvent(event)

		// регистрация клиента
		case client := <-h.register:
			h.clients[client] = struct{}{}
			client.send(map[string]string{
				"type":  "state",
				"state": h.fight.State().String(),
			})

		// отключение клиента
		case client := <-h.unregister:
			delete(h.clients, client)
		}
	}
}

// точка входа событий
func (h *Hub) Publish(event Event) {
	h.events <- event
}

// router доменных событий
func (h *Hub) handleEvent(event Event) {
	switch event.Type {
	case EventFightControl:
		h.handleFightControl(event.Data)
	case EventScore:
		h.handleScore(event.Data)
	case EventWarning:
		h.handleWarning(event.Data)
	case EventReset:
		h.handleReset()
	default:
		log.Printf("unknown event type: %v", event.Type)
	}
}

func (h *Hub) handleWarning(data json.RawMessage) {
	// нельзя штрафовать вне боя
	if h.fight.State() != fight.StateRunning {
		log.Println("warning ignored: fight not running")
		return
	}

	var payload WarningPayload
	if err := json.Unmarshal(data, &payload); err != nil {
		log.Printf("invalid warning payload: %v", err)
		return
	}

	event := events.WarningEvent{
		Fighter: events.Fighter(payload.Fighter),
		Time:    time.Now(),
	}

	// event sourcing
	h.eventLog = append(h.eventLog, event)

	penalty := h.warnings.Add(event.Fighter)

	if penalty {
		// штраф −1 балл
		h.scoreboard.Apply(events.ScoreEvent{
			Fighter: event.Fighter,
			Points:  -1,
			Time:    time.Now(),
		})

		log.Printf("PENALTY: fighter=%s -1 point", event.Fighter)
		h.broadcastScore()
	}

	red, blue := h.warnings.Count()
	log.Printf("WARNING: red=%d blue=%d", red, blue)

	h.broadcastWarnings()
}

func (h *Hub) handleFightControl(data json.RawMessage) {
	var evt FightControlEvent

	if err := json.Unmarshal(data, &evt); err != nil {
		log.Printf("invalid fight control payload: %v", err)
		return
	}

	var err error

	// FSM → Running
	switch evt.Action {
	case ActionStart:
		err = h.fight.Start()
		if err == nil {
			h.timer.Reset()
			h.timer.Start()
		}

	case ActionPause:
		err = h.fight.Pause()
		if err == nil {
			h.timer.Pause()
		}

	case ActionResume:
		err = h.fight.Resume()
		if err == nil {
			h.timer.Start()
		}

	// FSM → Stopped
	case ActionStop:
		err = h.fight.Stop()
		if err == nil {
			h.timer.Stop()
		}
	}

	if err != nil {
		log.Printf("fight action error: %v", err)
		return
	}

	log.Printf("fight state changed to %s", h.fight.State())
	h.broadcastState()
}

func (h *Hub) handleScore(data json.RawMessage) {
	// очки только в бою
	if h.fight.State() != fight.StateRunning {
		log.Println("score ignored: fight not running")
		return
	}

	var payload ScorePayload
	if err := json.Unmarshal(data, &payload); err != nil {
		log.Printf("invalid score payload: %v", err)
		return
	}

	event := events.ScoreEvent{
		Fighter: events.Fighter(payload.Fighter),
		Points:  payload.Points,
		Time:    time.Now(),
	}

	h.eventLog = append(h.eventLog, event)
	h.scoreboard.Apply(event)
	h.broadcastScore()
}

func (h *Hub) broadcastState() {
	for c := range h.clients {
		c.send(map[string]string{
			"type":  "state",
			"state": h.fight.State().String(),
		})
	}
}

func (h *Hub) broadcastScore() {
	red, blue := h.scoreboard.Score()

	for c := range h.clients {
		c.send(map[string]any{
			"type": "score_update",
			"red":  red,
			"blue": blue,
		})
	}
}

func (h *Hub) broadcastWarnings() {
	red, blue := h.warnings.Count()

	for c := range h.clients {
		c.send(map[string]any{
			"type": "warnings",
			"red":  red,
			"blue": blue,
		})
	}
}

func (h *Hub) broadcastTimer(seconds int) {
	for c := range h.clients {
		c.send(map[string]any{
			"type":         "timer",
			"seconds_left": seconds,
		})
	}
}
