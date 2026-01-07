package fight

import "time"

type Timer struct {
	duration   time.Duration
	remaining  time.Duration
	ticker     *time.Ticker
	stop       chan struct{}
	onTick     func(time.Duration)
	onFinished func()
}

func NewTimer(d time.Duration) *Timer {
	return &Timer{
		duration:  d,
		remaining: d,
	}
}

func (t *Timer) SetDuration(d time.Duration) {
	t.duration = d
	t.remaining = d
}

func (t *Timer) Start() {
	if t.ticker != nil {
		return
	}

	t.stop = make(chan struct{})
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
					t.stopInternal()
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

func (t *Timer) stopInternal() {
	if t.ticker != nil {
		t.ticker.Stop()
		t.ticker = nil
	}
	close(t.stop)
}

func (t *Timer) Stop() {
	if t.ticker == nil {
		return
	}
	t.stopInternal()
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
