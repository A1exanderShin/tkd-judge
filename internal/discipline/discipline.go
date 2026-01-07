package discipline

type RealtimeEvent struct {
	Type string
	Data any
}

type Discipline interface {
	HandleEvent(event any) error
	Snapshot() map[string]any
	Reset()

	Realtime() <-chan RealtimeEvent
}
