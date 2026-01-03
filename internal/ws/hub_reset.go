package ws

import (
	"log"

	"tkd-judge/internal/fight"
)

func (h *Hub) handleReset() {
	state := h.fight.State()

	if state != fight.StateFinished && state != fight.StatePaused {
		log.Println("reset ignored: invalid state")
		return
	}

	h.timer.Stop()

	h.fight = fight.NewFight()
	h.timer.Reset()
	h.scoreboard = fight.NewScoreboard()
	h.warnings = fight.NewWarningCounter(h.cfg.WarningsForPenalty)
	h.eventLog = make([]any, 0)

	log.Println("FIGHT RESET")

	h.broadcastState()
	h.broadcastScore()
	h.broadcastWarnings()
	h.broadcastTimer(h.timer.RemainingSeconds())
}
