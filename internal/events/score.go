package events

import "time"

type Fighter string

const (
	FighterRed  Fighter = "red"
	FighterBlue Fighter = "blue"
)

type ScoreEvent struct {
	JudgeID int
	Fighter Fighter
	Points  int
	Time    time.Time
}
