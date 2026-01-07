package discipline

type PatternEventType string

const (
	EventPatternSelectExercise PatternEventType = "PATTERN_SELECT_EXERCISE"
	EventPatternJudgeScore     PatternEventType = "PATTERN_JUDGE_SCORE"
	EventPatternNextStage      PatternEventType = "PATTERN_NEXT_STAGE"
	EventPatternFinish         PatternEventType = "PATTERN_FINISH"
	EventPatternReset          PatternEventType = "PATTERN_RESET"
)

type PatternEvent struct {
	Type PatternEventType

	Exercise  string
	JudgeID   int
	Criterion string
	Score     int
}
