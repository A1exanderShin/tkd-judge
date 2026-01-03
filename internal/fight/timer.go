package fight

import "time"

type Timer struct {
	duration   time.Duration
	remaining  time.Duration
	ticker     *time.Ticker
	stop       chan struct{}
	onTick     func(remaining time.Duration)
	onFinished func()
}

func NewTimer(duration time.Duration) *Timer {
	return &Timer{
		duration:  duration,
		remaining: duration,
		stop:      make(chan struct{}),
	}
}

func (t *Timer) Start() {
	if t.ticker != nil {
		return
	}

	t.ticker = time.NewTicker(time.Second)

	go func() {
		for {
			select {
			case <-t.ticker.C:
				t.remaining -= time.Second

				if t.onTick != nil {
					t.onTick(t.remaining)
				}

				if t.remaining <= 0 {
					t.Stop()
					if t.onFinished != nil {
						t.onFinished()
					}
					return
				}

			case <-t.stop:
				return
			}
		}
	}()
}

func (t *Timer) Pause() {
	t.Stop()
}

func (t *Timer) Stop() {
	if t.ticker != nil {
		t.ticker.Stop()
		t.ticker = nil
	}
}

func (t *Timer) Reset() {
	t.Stop()
	t.remaining = t.duration
}

func (t *Timer) RemainingSeconds() int {
	return int(t.remaining.Seconds())
}

func (t *Timer) OnTick(fn func(time.Duration)) {
	t.onTick = fn
}

func (t *Timer) OnFinished(fn func()) {
	t.onFinished = fn
}
