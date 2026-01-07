package discipline

import "tkd-judge/internal/events"

type FightEventType string

const (
	EventFightStart    FightEventType = "FIGHT_START"
	EventFightPause    FightEventType = "FIGHT_PAUSE"
	EventFightStop     FightEventType = "FIGHT_STOP"
	EventFightReset    FightEventType = "FIGHT_RESET"
	EventFightScore    FightEventType = "FIGHT_SCORE"
	EventFightSettings FightEventType = "FIGHT_SETTINGS"
	EventFightWarning  FightEventType = "FIGHT_WARNING"
)

type FightEvent struct {
	Type FightEventType

	// ===== score =====
	JudgeID int
	Fighter events.Fighter
	Points  int

	// ===== settings =====
	Rounds        int
	RoundDuration int // seconds
	BreakDuration int // seconds
}
