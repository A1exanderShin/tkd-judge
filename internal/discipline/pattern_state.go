package discipline

import "errors"

type PatternState int

const (
	StateIdle PatternState = iota
	StateJudging
	StateResult
	StateFinished
)

func (s PatternState) String() string {
	switch s {
	case StateIdle:
		return "IDLE"
	case StateJudging:
		return "JUDGING"
	case StateResult:
		return "RESULT"
	case StateFinished:
		return "FINISHED"
	default:
		return "UNKNOWN"
	}
}

var (
	ErrPatternFinished = errors.New("pattern already finished")
)
