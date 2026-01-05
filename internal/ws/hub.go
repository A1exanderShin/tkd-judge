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

type hubEvent struct {
	event  Event
	client *Client
}

type judgeScore struct {
	red  int
	blue int
}

type Hub struct {
	cfg config.Config

	fight *fight.Fight

	roundTimer *fight.Timer
	breakTimer *fight.Timer

	roundDuration time.Duration
	breakDuration time.Duration

	scoreboard *fight.Scoreboard
	warnings   *fight.WarningCounter

	judges      map[int]*judges.Judge
	judgeScores map[int]judgeScore

	events chan hubEvent

	clients    map[*Client]struct{}
	mainJudge  *Client
	sideJudges map[int]*Client

	register   chan *Client
	unregister chan *Client
}

/* ================= CONSTRUCTOR ================= */

func NewHub() *Hub {
	cfg := config.Default()

	h := &Hub{
		cfg:           cfg,
		fight:         fight.NewFight(),
		roundDuration: cfg.RoundDuration,
		breakDuration: 30 * time.Second,

		scoreboard:  fight.NewScoreboard(),
		warnings:    fight.NewWarningCounter(cfg.WarningsForPenalty),
		judges:      make(map[int]*judges.Judge),
		judgeScores: make(map[int]judgeScore),

		events:     make(chan hubEvent, 16),
		clients:    make(map[*Client]struct{}),
		sideJudges: make(map[int]*Client),

		register:   make(chan *Client),
		unregister: make(chan *Client),
	}

	for i := 1; i <= cfg.JudgesCount; i++ {
		h.judges[i] = judges.NewJudge(i, cfg.AntiClick)
	}

	// ===== timers =====
	h.roundTimer = fight.NewTimer(h.roundDuration)
	h.breakTimer = fight.NewTimer(h.breakDuration)

	h.roundTimer.OnTick(func(rem time.Duration) {
		h.broadcastTimer(int(rem.Seconds()), "round")
	})
	h.breakTimer.OnTick(func(rem time.Duration) {
		h.broadcastTimer(int(rem.Seconds()), "break")
	})

	h.roundTimer.OnFinished(h.onRoundFinished)
	h.breakTimer.OnFinished(h.onBreakFinished)

	return h
}

/* ================= CORE LOOP ================= */

func (h *Hub) Publish(e Event, c *Client) {
	h.events <- hubEvent{event: e, client: c}
}

func (h *Hub) Run() {
	for {
		select {

		case ev := <-h.events:
			h.handleEvent(ev.event, ev.client)

		case c := <-h.register:
			if c.role == RoleMainJudge {
				if h.mainJudge != nil {
					c.close()
					continue
				}
				h.mainJudge = c
				log.Println("MAIN JUDGE CONNECTED")
			}

			if c.role == RoleSideJudge {
				h.sideJudges[c.judgeID] = c
				log.Printf("SIDE JUDGE %d CONNECTED", c.judgeID)
			}

			h.clients[c] = struct{}{}
			h.sendFullState(c)

		case c := <-h.unregister:
			delete(h.clients, c)

			if c == h.mainJudge {
				h.mainJudge = nil
			}
			if c.role == RoleSideJudge {
				delete(h.sideJudges, c.judgeID)
			}
		}
	}
}

/* ================= EVENT ROUTER ================= */

func (h *Hub) handleEvent(e Event, c *Client) {
	switch e.Type {

	case EventFightControl:
		if c.role == RoleMainJudge {
			h.handleFightControl(e.Data)
		}

	case EventFightSettings:
		if c.role == RoleMainJudge {
			h.handleFightSettings(e.Data)
		}

	case EventScore:
		if c.role == RoleSideJudge {
			h.handleScore(e.Data, c)
		}

	case EventReset:
		if c.role == RoleMainJudge {
			h.handleReset()
		}
	}
}

/* ================= FIGHT CONTROL ================= */

func (h *Hub) handleFightControl(data json.RawMessage) {
	var evt FightControlEvent
	if json.Unmarshal(data, &evt) != nil {
		return
	}

	switch evt.Action {

	case ActionStart:
		h.fight.Start()
		h.roundTimer.Reset()
		h.roundTimer.Start()

	case ActionPause:
		h.roundTimer.Stop()
		h.breakTimer.Stop()
		h.fight.Pause()

	case ActionStop:
		h.roundTimer.Stop()
		h.breakTimer.Stop()
		h.fight.Stop()
	}

	h.broadcastState()
}

/* ================= SETTINGS ================= */

func (h *Hub) handleFightSettings(data json.RawMessage) {
	if h.fight.State() == fight.StateRunning {
		return
	}

	var p FightSettingsPayload
	if json.Unmarshal(data, &p) != nil {
		return
	}

	if p.RoundDuration > 0 {
		h.roundDuration = time.Duration(p.RoundDuration) * time.Second
		h.roundTimer.SetDuration(h.roundDuration)
	}

	if p.BreakDuration > 0 {
		h.breakDuration = time.Duration(p.BreakDuration) * time.Second
		h.breakTimer.SetDuration(h.breakDuration)
	}

	if p.Rounds > 0 {
		h.fight.SetRounds(p.Rounds)
	}

	h.broadcastSettings()
}

/* ================= SCORE ================= */

func (h *Hub) handleScore(data json.RawMessage, c *Client) {
	if h.fight.State() != fight.StateRunning {
		return
	}

	var p ScorePayload
	if json.Unmarshal(data, &p) != nil {
		return
	}

	if h.judges[c.judgeID].CanScore(time.Now()) != nil {
		return
	}

	ev := events.ScoreEvent{
		JudgeID: c.judgeID,
		Fighter: events.Fighter(p.Fighter),
		Points:  p.Points,
		Time:    time.Now(),
	}

	h.scoreboard.Apply(ev)

	js := h.judgeScores[c.judgeID]
	if p.Fighter == "red" {
		js.red += p.Points
	} else {
		js.blue += p.Points
	}
	h.judgeScores[c.judgeID] = js

	h.broadcastScore()
	h.broadcastJudgeScores()
}

/* ================= TIMERS ================= */

func (h *Hub) onRoundFinished() {
	if h.fight.CurrentRound() < h.fight.TotalRounds() {
		h.fight.SetState(fight.StateBreak)
		h.breakTimer.Reset()
		h.breakTimer.Start()
	} else {
		h.fight.Stop()
	}
	h.broadcastState()
}

func (h *Hub) onBreakFinished() {
	// увеличиваем номер раунда
	h.fight.NextRound()

	// переводим в ожидание старта
	h.fight.SetState(fight.StatePaused)

	// 🔥 ВАЖНО: шлём обновлённые настройки боя
	h.broadcastSettings()

	// состояние
	h.broadcastState()
}

/* ================= RESET ================= */

func (h *Hub) handleReset() {
	h.roundTimer.Stop()
	h.breakTimer.Stop()

	h.fight = fight.NewFight()
	h.scoreboard = fight.NewScoreboard()
	h.warnings = fight.NewWarningCounter(h.cfg.WarningsForPenalty)
	h.judgeScores = make(map[int]judgeScore)

	h.broadcastState()
	h.broadcastScore()
	h.broadcastJudgeScores()
}

/* ================= BROADCAST ================= */

func (h *Hub) sendFullState(c *Client) {
	c.send(map[string]any{
		"type":  "state",
		"state": h.fight.State().String(),
	})

	red, blue := h.scoreboard.Score()
	c.send(map[string]any{
		"type": "score_update",
		"red":  red,
		"blue": blue,
	})

	h.broadcastSettingsTo(c)

	if c.role == RoleSideJudge {
		c.send(map[string]any{
			"type":    "judge_id",
			"judgeID": c.judgeID,
		})
	}

	h.broadcastJudgeScoresTo(c)
}

func (h *Hub) broadcastState() {
	for c := range h.clients {
		c.send(map[string]any{
			"type":  "state",
			"state": h.fight.State().String(),
		})
	}
}

func (h *Hub) broadcastScore() {
	for c := range h.clients {
		red, blue := h.scoreboard.Score()
		c.send(map[string]any{
			"type": "score_update",
			"red":  red,
			"blue": blue,
		})
	}
}

func (h *Hub) broadcastJudgeScores() {
	for c := range h.clients {
		h.broadcastJudgeScoresTo(c)
	}
}

func (h *Hub) broadcastJudgeScoresTo(c *Client) {
	type row struct {
		ID   int `json:"id"`
		Red  int `json:"red"`
		Blue int `json:"blue"`
	}

	out := []row{}
	for i := 1; i <= h.cfg.JudgesCount; i++ {
		js := h.judgeScores[i]
		out = append(out, row{ID: i, Red: js.red, Blue: js.blue})
	}

	c.send(map[string]any{
		"type":   "judge_scores",
		"scores": out,
	})
}

func (h *Hub) broadcastSettings() {
	for c := range h.clients {
		h.broadcastSettingsTo(c)
	}
}

func (h *Hub) broadcastSettingsTo(c *Client) {
	c.send(map[string]any{
		"type":           "fight_settings",
		"round_duration": int(h.roundDuration.Seconds()),
		"break_duration": int(h.breakDuration.Seconds()),
		"round":          h.fight.CurrentRound(),
		"rounds_total":   h.fight.TotalRounds(),
	})
}

func (h *Hub) broadcastTimer(seconds int, mode string) {
	for c := range h.clients {
		c.send(map[string]any{
			"type":         "timer",
			"seconds_left": seconds,
			"mode":         mode, // round | break
		})
	}
}

/* ================= HELPERS ================= */

func (h *Hub) nextFreeJudgeID() int {
	for i := 1; i <= h.cfg.JudgesCount; i++ {
		if _, ok := h.sideJudges[i]; !ok {
			return i
		}
	}
	return 0
}
