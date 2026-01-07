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

	// score
	Fighter string `json:"fighter,omitempty"`
	Points  int    `json:"points,omitempty"`

	// settings
	Rounds        int `json:"rounds,omitempty"`
	RoundDuration int `json:"round_duration,omitempty"`
	BreakDuration int `json:"break_duration,omitempty"`
}

/* ================= PARSER ================= */

func ParseMessage(raw json.RawMessage, c *Client) (any, error) {
	var msg IncomingMessage
	if err := json.Unmarshal(raw, &msg); err != nil {
		return nil, err
	}

	switch msg.Type {

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

		f := events.Fighter(msg.Fighter)
		if f != events.FighterRed && f != events.FighterBlue {
			return nil, nil
		}

		return discipline.FightEvent{
			Type:    discipline.EventFightScore,
			JudgeID: c.judgeID,
			Fighter: f,
			Points:  msg.Points,
		}, nil

	// 🔥🔥🔥 ВОТ ОН — КЛЮЧЕВОЙ КЕЙС
	case "FIGHT_WARNING":
		f := events.Fighter(msg.Fighter)
		if f != events.FighterRed && f != events.FighterBlue {
			return nil, nil
		}

		return discipline.FightEvent{
			Type:    discipline.EventFightWarning,
			Fighter: f,
		}, nil
	}

	return nil, nil
}
