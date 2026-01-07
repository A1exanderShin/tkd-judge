package discipline

/* ================= STRUCT ================= */

type PatternDiscipline struct {
	pattern        *Pattern
	expectedJudges int
	realtime       chan RealtimeEvent
}

/* ================= CONSTRUCTOR ================= */

func NewPatternDiscipline(criteria []string, judges int) *PatternDiscipline {
	return &PatternDiscipline{
		pattern:        NewPattern(criteria),
		expectedJudges: judges,
		realtime:       make(chan RealtimeEvent, 8),
	}
}

/* ================= DISCIPLINE INTERFACE ================= */

func (pd *PatternDiscipline) HandleEvent(event any) error {
	e, ok := event.(PatternEvent)
	if !ok {
		return nil
	}

	// 🔥 RESET вне FSM
	if e.Type == EventPatternReset {
		pd.Reset()
		return nil
	}

	switch pd.pattern.state {

	case StateIdle:
		return pd.handleIdle(e)

	case StateJudging:
		return pd.handleJudging(e)

	case StateResult:
		return pd.handleResult(e)

	case StateFinished:
		return ErrPatternFinished
	}

	return nil
}

func (pd *PatternDiscipline) Snapshot() map[string]any {
	return map[string]any{
		"type":     "pattern",
		"state":    pd.pattern.state.String(),
		"stage":    pd.pattern.stage,
		"exercise": pd.pattern.currentExercise,
		"total":    pd.pattern.CalculateTotal(),
	}
}

func (pd *PatternDiscipline) Realtime() <-chan RealtimeEvent {
	return pd.realtime
}

func (pd *PatternDiscipline) Reset() {
	pd.pattern = NewPattern(pd.pattern.criteria)
	pd.emit("reset")
}

/* ================= FSM ================= */

func (pd *PatternDiscipline) handleIdle(e PatternEvent) error {
	if e.Type != EventPatternSelectExercise {
		return nil
	}

	if e.Exercise == "" {
		return nil
	}

	pd.pattern.SetExercise(e.Exercise)
	pd.pattern.state = StateJudging

	pd.emit("exercise_selected")
	return nil
}

func (pd *PatternDiscipline) handleJudging(e PatternEvent) error {
	if e.Type != EventPatternJudgeScore {
		return nil
	}

	if e.JudgeID <= 0 {
		return nil // защита от кривого UI
	}

	pd.pattern.AddScore(e.JudgeID, e.Criterion, e.Score)

	if pd.pattern.AllJudgesVoted(pd.expectedJudges) {
		pd.pattern.state = StateResult
		pd.emit("result_ready")
	}

	return nil
}

func (pd *PatternDiscipline) handleResult(e PatternEvent) error {
	switch e.Type {

	case EventPatternNextStage:
		pd.pattern.NextStage()
		pd.pattern.state = StateIdle
		pd.emit("next_stage")

	case EventPatternFinish:
		pd.pattern.Finish()
		pd.emit("finished")
	}

	return nil
}

/* ================= REALTIME ================= */

func (pd *PatternDiscipline) emit(kind string) {
	select {
	case pd.realtime <- RealtimeEvent{
		Type: kind,
		Data: pd.Snapshot(),
	}:
	default:
		// не блокируем FSM
	}
}
