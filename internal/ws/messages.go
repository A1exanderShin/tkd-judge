package ws

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

// Универсальное событие, которое попадает в hub
type Event struct {
	Type EventType
	Data any
}

// Payload управления боем
type FightControlEvent struct {
	Action FightAction
}
