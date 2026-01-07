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
	if antiClick < 0 {
		antiClick = 0
	}

	return &Judge{
		ID:        id,
		antiClick: antiClick,
	}
}

func (j *Judge) CanScore(now time.Time) error {
	if j.antiClick == 0 {
		j.lastClickAt = now
		return nil
	}

	if !j.lastClickAt.IsZero() && now.Sub(j.lastClickAt) < j.antiClick {
		return ErrTooFast
	}

	j.lastClickAt = now
	return nil
}

func (j *Judge) Reset() {
	j.lastClickAt = time.Time{}
}
