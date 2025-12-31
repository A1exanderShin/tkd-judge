package fight

type Fight struct {
	state State
}

// NewFight создаёт новый бой в состоянии IDLE
func NewFight() *Fight {
	// TODO: вернуть бой со StateIdle
	return &Fight{state: StateIdle}
}

// Start переводит бой из IDLE в RUNNING
func (f *Fight) Start() error {
	// TODO:
	// - если состояние не IDLE → ошибка
	// - если ок → StateRunning
	if f.state != StateIdle {
		return ErrInvalidTransition
	}
	f.state = StateRunning
	return nil
}

// Pause переводит бой из RUNNING в PAUSED
func (f *Fight) Pause() error {
	// TODO:
	// - если состояние не RUNNING → ошибка
	// - если ок → StatePaused
	if f.state != StateRunning {
		return ErrInvalidTransition
	}
	f.state = StatePaused
	return nil
}

// Resume переводит бой из PAUSED в RUNNING
func (f *Fight) Resume() error {
	// TODO:
	// - если состояние не PAUSED → ошибка
	// - если ок → StateRunning
	if f.state != StatePaused {
		return ErrInvalidTransition
	}
	f.state = StateRunning
	return nil
}

// Stop переводит бой из RUNNING в FINISHED
func (f *Fight) Stop() error {
	// TODO:
	// - если состояние уже FINISHED → ErrFightFinished
	// - если состояние не RUNNING → ErrInvalidTransition
	// - если ок → StateFinished
	if f.state == StateFinished {
		return ErrFightFinished
	}

	if f.state != StateRunning {
		return ErrInvalidTransition
	}
	f.state = StateFinished
	return nil
}

// State возвращает текущее состояние боя
func (f *Fight) State() State {
	// TODO: вернуть состояние
	return f.state
}
