package discipline

import "tkd-judge/internal/events"

type FightEventType string

const (
	EventFightStart FightEventType = "FIGHT_START"
	EventFightPause FightEventType = "FIGHT_PAUSE"
	EventFightStop  FightEventType = "FIGHT_STOP"
	EventFightScore FightEventType = "FIGHT_SCORE"
	EventFightReset FightEventType = "FIGHT_RESET"
)

type FightEvent struct {
	Type FightEventType

	JudgeID int
	Fighter events.Fighter
	Points  int
}
