package ws

import (
	"encoding/json"

	"tkd-judge/internal/discipline"
	"tkd-judge/internal/events"
)

/*
WS EVENT FORMAT (JSON):

{
  "type": "FIGHT_START"
}

{
  "type": "FIGHT_SCORE",
  "fighter": "red",
  "points": 1
}

{
  "type": "FIGHT_SETTINGS",
  "rounds": 3,
  "round_duration": 20,
  "break_duration": 20
}
*/

type IncomingMessage struct {
	Type string `json:"type"`

	// fight
	Fighter string `json:"fighter,omitempty"`
	Points  int    `json:"points,omitempty"`

	// pattern
	Exercise  string `json:"exercise,omitempty"`
	Criterion string `json:"criterion,omitempty"`
	Score     int    `json:"score,omitempty"`
}

/* ================= PARSER ================= */

func ParseMessage(raw json.RawMessage, c *Client) (any, error) {
	var msg IncomingMessage
	if err := json.Unmarshal(raw, &msg); err != nil {
		return nil, err
	}

	switch msg.Type {

	// ================= FIGHT =================

	case "FIGHT_START":
		return discipline.FightEvent{Type: discipline.EventFightStart}, nil

	case "FIGHT_PAUSE":
		return discipline.FightEvent{Type: discipline.EventFightPause}, nil

	case "FIGHT_STOP":
		return discipline.FightEvent{Type: discipline.EventFightStop}, nil

	case "FIGHT_RESET":
		return discipline.FightEvent{Type: discipline.EventFightReset}, nil

	case "FIGHT_SCORE":
		if c.judgeID == 0 {
			return nil, nil
		}
		return discipline.FightEvent{
			Type:    discipline.EventFightScore,
			JudgeID: c.judgeID,
			Fighter: events.Fighter(msg.Fighter),
			Points:  msg.Points,
		}, nil

	// ================= PATTERN =================

	case "PATTERN_SELECT_EXERCISE":
		return discipline.PatternEvent{
			Type:     discipline.EventPatternSelectExercise,
			Exercise: msg.Exercise,
		}, nil

	case "PATTERN_JUDGE_SCORE":
		if c.judgeID == 0 {
			return nil, nil
		}
		return discipline.PatternEvent{
			Type:      discipline.EventPatternJudgeScore,
			JudgeID:   c.judgeID,
			Criterion: msg.Criterion,
			Score:     msg.Score,
		}, nil

	case "PATTERN_NEXT_STAGE":
		return discipline.PatternEvent{
			Type: discipline.EventPatternNextStage,
		}, nil

	case "PATTERN_FINISH":
		return discipline.PatternEvent{
			Type: discipline.EventPatternFinish,
		}, nil

	case "SWITCH_DISCIPLINE":
		return SystemEvent{
			Type: EventSwitchDiscipline,
		}, nil
	}

	return nil, nil
}
