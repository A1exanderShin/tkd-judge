package fight

import "errors"

type State int

const (
	StateIdle State = iota
	StateRunning
	StatePaused
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
	case StateFinished:
		return "FINISHED"
	default:
		return "UNKNOWN"
	}
}

var (
	// ErrInvalidTransition — попытка недопустимого перехода FSM
	ErrInvalidTransition = errors.New("invalid fight state transition")

	// ErrFightFinished — любые действия после завершения боя
	ErrFightFinished = errors.New("fight already finished")
)
