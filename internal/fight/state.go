package fight

import "errors"

type State int

const (
	StateIdle State = iota
	StateRunning
	StatePaused
	StateBreak
	StateFinished
)

// String возвращает человекочитаемое название состояния.
func (s State) String() string {
	switch s {
	case StateIdle:
		return "IDLE"
	case StateRunning:
		return "RUNNING"
	case StatePaused:
		return "PAUSED"
	case StateBreak:
		return "BREAK"
	case StateFinished:
		return "FINISHED"
	default:
		return "UNKNOWN"
	}
}

var (
	// ErrFightFinished — любые действия после завершения боя
	ErrFightFinished = errors.New("fight already finished")
)
