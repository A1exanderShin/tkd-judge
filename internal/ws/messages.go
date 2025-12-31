package ws

import "encoding/json"

type EventType string

const (
	EventFightControl EventType = "fight_control"
)

type FightAction string

const (
	ActionStart  FightAction = "start"
	ActionPause  FightAction = "pause"
	ActionResume FightAction = "resume"
	ActionStop   FightAction = "stop"
)

type Event struct {
	Type EventType       `json:"Type"`
	Data json.RawMessage `json:"Data"`
}

type FightControlEvent struct {
	Action FightAction `json:"Action"`
}
