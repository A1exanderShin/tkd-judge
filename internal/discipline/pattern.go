package discipline

type Pattern struct {
	state           PatternState
	stage           int
	currentExercise string
	criteria        []string
	scores          map[int]map[string]int
}

func NewPattern() *Pattern {
	return &Pattern{
		state:           StateIdle,
		stage:           0,
		currentExercise: "",
		criteria:        make([]string, 0),
		scores:          make(map[int]map[string]int),
	}
}

func (p *Pattern) SetCurrentExercise(exercise string) {
	p.currentExercise = exercise
	p.scores = make(map[int]map[string]int)
}

func (p *Pattern) AddScore(judgeID int, criterion string, score int) {
	if _, ok := p.scores[judgeID]; !ok {
		p.scores[judgeID] = make(map[string]int)
	}

	p.scores[judgeID][criterion] = score
}

func (p *Pattern) AllJudgesVoted(expectedJudges int) bool {
	if len(p.scores) < expectedJudges {
		return false
	}

	for _, criteriaScores := range p.scores {
		if len(criteriaScores) < len(p.criteria) {
			return false
		}
	}

	return true
}

func (p *Pattern) CalculateResult() int {
	total := 0

	for _, criteriaScores := range p.scores {
		for _, score := range criteriaScores {
			total += score
		}
	}

	return total
}

func (p *Pattern) NextStage() {
	p.stage++
	p.state = StateIdle
	p.currentExercise = ""
	p.scores = make(map[int]map[string]int)
}

func (p *Pattern) Finish() {
	p.state = StateFinished
}

func (p *Pattern) State() PatternState {
	return p.state
}

func (p *Pattern) Stage() int {
	return p.stage
}

func (p *Pattern) Reset() {
	p.state = StateIdle
	p.stage = 0
	p.currentExercise = ""
	p.scores = make(map[int]map[string]int)
}

func (p *Pattern) Snapshot() map[string]any {
	return map[string]any{
		"state":    p.state.String(),
		"stage":    p.stage,
		"exercise": p.currentExercise,
		"scores":   p.scores,
		"criteria": p.criteria,
	}
}
