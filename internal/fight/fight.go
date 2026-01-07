package fight

type Fight struct {
	state        State
	currentRound int
	totalRounds  int
}

func NewFight() *Fight {
	return &Fight{
		state:        StateIdle,
		currentRound: 1,
		totalRounds:  1,
	}
}

func (f *Fight) SetRounds(total int) {
	if total < 1 {
		total = 1
	}
	f.totalRounds = total
	f.currentRound = 1
}

func (f *Fight) Start() {
	// разрешаем старт только из IDLE или PAUSED
	if f.state != StateIdle && f.state != StatePaused {
		return
	}
	f.state = StateRunning
}

func (f *Fight) Pause() {
	if f.state == StateRunning {
		f.state = StatePaused
	}
}

func (f *Fight) Stop() {
	if f.state == StateFinished {
		return
	}
	f.state = StateFinished
}

func (f *Fight) SetState(s State) {
	f.state = s
}

func (f *Fight) State() State {
	return f.state
}

func (f *Fight) CurrentRound() int {
	return f.currentRound
}

func (f *Fight) TotalRounds() int {
	return f.totalRounds
}

func (f *Fight) NextRound() {
	if f.currentRound < f.totalRounds {
		f.currentRound++
	}
}
