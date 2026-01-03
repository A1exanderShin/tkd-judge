package judges

import (
	"errors"
	"time"
)

var ErrTooFast = errors.New("click too fast")

type Judge struct {
	ID          int
	lastClickAt time.Time
	antiClick   time.Duration
}

func NewJudge(id int, antiClick time.Duration) *Judge {
	return &Judge{
		ID:        id,
		antiClick: antiClick,
	}
}

func (j *Judge) CanScore(now time.Time) error {
	if !j.lastClickAt.IsZero() && now.Sub(j.lastClickAt) < j.antiClick {
		return ErrTooFast
	}

	j.lastClickAt = now
	return nil
}
