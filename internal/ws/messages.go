package ws

import "encoding/json"

type EventType string

const (
	EventFightControl  EventType = "fight_control"
	EventFightSettings EventType = "fight_settings"
	EventScore         EventType = "score"
	EventWarning       EventType = "warning"
	EventReset         EventType = "reset"
)

type Event struct {
	Type EventType       `json:"Type"`
	Data json.RawMessage `json:"Data"`
}

/* ===== fight control ===== */

type FightAction string

const (
	ActionStart  FightAction = "start"
	ActionPause  FightAction = "pause"
	ActionResume FightAction = "resume"
	ActionStop   FightAction = "stop"
)

type FightControlEvent struct {
	Action FightAction `json:"Action"`
}

/* ===== settings ===== */

type FightSettingsPayload struct {
	RoundDuration int `json:"round_duration"` // sec
	Rounds        int `json:"rounds"`
	BreakDuration int `json:"break_duration"` // sec
}

/* ===== score ===== */

type ScorePayload struct {
	Fighter string `json:"Fighter"`
	Points  int    `json:"Points"`
}

/* ===== warning ===== */

type WarningPayload struct {
	Fighter string `json:"Fighter"`
}
