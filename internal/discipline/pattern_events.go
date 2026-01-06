package discipline

type PatternEventType string

const (
	EventSelectExercise PatternEventType = "SELECT_EXERCISE"
	EventJudgeScore     PatternEventType = "JUDGE_SCORE"
	EventNextStage      PatternEventType = "NEXT_STAGE"
	EventFinish         PatternEventType = "FINISH"
)

type PatternEvent struct {
	Type PatternEventType

	// payload
	Exercise  string
	JudgeID   int
	Criterion string
	Score     int
}

func (p *Pattern) HandleEvent(event PatternEvent, expectedJudges int) error {
	switch p.state {

	case StateIdle:
		return p.handleIdle(event)

	case StateJudging:
		return p.handleJudging(event, expectedJudges)

	case StateResult:
		return p.handleResult(event)

	case StateFinished:
		return ErrPatternFinished
	}

	return nil
}

func (p *Pattern) handleIdle(event PatternEvent) error {
	switch event.Type {

	case EventSelectExercise:
		p.SetCurrentExercise(event.Exercise)
		p.state = StateJudging
		return nil

	default:
		return nil // игнорируем
	}
}

func (p *Pattern) handleJudging(event PatternEvent, expectedJudges int) error {
	switch event.Type {

	case EventJudgeScore:
		p.AddScore(event.JudgeID, event.Criterion, event.Score)

		if p.AllJudgesVoted(expectedJudges) {
			p.state = StateResult
		}
		return nil

	default:
		return nil
	}
}

func (p *Pattern) handleResult(event PatternEvent) error {
	switch event.Type {

	case EventNextStage:
		p.NextStage()
		p.state = StateIdle
		return nil

	case EventFinish:
		p.Finish()
		return nil

	default:
		return nil
	}
}

func (p *Pattern) HandleAnyEvent(event any) error {
	e, ok := event.(PatternEvent)
	if !ok {
		return nil
	}

	// expectedJudges пока передаётся снаружи
	return p.HandleEvent(e /* expectedJudges */, 5)
}
