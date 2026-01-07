package ws

type SystemEventType string

const (
	EventSwitchDiscipline SystemEventType = "SWITCH_DISCIPLINE"
)

type SystemEvent struct {
	Type SystemEventType
}
