package discipline

type Pattern struct {
	state           PatternState
	stage           int
	currentExercise string
	criteria        []string
	scores          map[int]map[string]int
}

func NewPattern(criteria []string) *Pattern {
	return &Pattern{
		state:    StateIdle,
		stage:    1,
		criteria: criteria,
		scores:   make(map[int]map[string]int),
	}
}

/* ===== getters ===== */

func (p *Pattern) State() PatternState {
	return p.state
}

func (p *Pattern) Stage() int {
	return p.stage
}

func (p *Pattern) Exercise() string {
	return p.currentExercise
}

/* ===== domain logic ===== */

func (p *Pattern) SetExercise(ex string) {
	p.currentExercise = ex
	p.state = StateJudging
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
	for _, m := range p.scores {
		if len(m) < len(p.criteria) {
			return false
		}
	}
	return true
}

func (p *Pattern) CalculateTotal() int {
	total := 0
	for _, m := range p.scores {
		for _, v := range m {
			total += v
		}
	}
	return total
}

func (p *Pattern) NextStage() {
	p.stage++
	p.currentExercise = ""
	p.scores = make(map[int]map[string]int)
}

func (p *Pattern) Finish() {
	p.state = StateFinished
}

func (p *Pattern) ResetStage() {
	p.scores = make(map[int]map[string]int)
}
