package ws

import (
	"time"

	"tkd-judge/internal/events"
	"tkd-judge/internal/protocol"
)

func (h *Hub) BuildProtocol() protocol.Protocol {
	var p protocol.Protocol

	p.GeneratedAt = time.Now()
	p.State = h.fight.State().String()

	red, blue := h.scoreboard.Score()
	p.Score.Red = red
	p.Score.Blue = blue

	wr, wb := h.warnings.Count()
	p.Warnings.Red = wr
	p.Warnings.Blue = wb

	for _, e := range h.eventLog {
		switch v := e.(type) {

		case events.ScoreEvent:
			p.Events = append(p.Events, protocol.ScoreLog{
				Type:    "score",
				JudgeID: v.JudgeID,
				Fighter: v.Fighter,
				Points:  v.Points,
				Time:    v.Time,
			})

		case events.WarningEvent:
			p.Events = append(p.Events, protocol.WarningLog{
				Type:    "warning",
				Fighter: v.Fighter,
				Time:    v.Time,
			})
		}
	}

	return p
}
