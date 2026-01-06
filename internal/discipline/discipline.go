package discipline

type Discipline interface {
	HandleEvent(event any) error
	Snapshot() map[string]any
	Reset()
}
