package discipline

import (
	"time"

	"tkd-judge/internal/config"
	"tkd-judge/internal/events"
	"tkd-judge/internal/fight"
)

type FightDiscipline struct {
	cfg config.Config

	fight *fight.Fight

	scoreboard *fight.Scoreboard
	warnings   *fight.WarningCounter

	roundTimer *fight.Timer
	breakTimer *fight.Timer

	realtime chan RealtimeEvent
}

/* ================= CONSTRUCTOR ================= */

func NewFightDiscipline() *FightDiscipline {
	cfg := config.Default()

	fd := &FightDiscipline{
		cfg:        cfg,
		fight:      fight.NewFight(),
		scoreboard: fight.NewScoreboard(),
		warnings:   fight.NewWarningCounter(cfg.WarningsForPenalty),
		realtime:   make(chan RealtimeEvent, 8),
	}

	fd.roundTimer = fight.NewTimer(cfg.RoundDuration)
	fd.breakTimer = fight.NewTimer(30 * time.Second)

	// realtime timer events
	fd.roundTimer.OnTick(func(rem time.Duration) {
		fd.emitTimer("round", rem)
	})
	fd.breakTimer.OnTick(func(rem time.Duration) {
		fd.emitTimer("break", rem)
	})

	fd.roundTimer.OnFinished(fd.onRoundFinished)
	fd.breakTimer.OnFinished(fd.onBreakFinished)

	return fd
}

/* ================= DISCIPLINE INTERFACE ================= */

func (fd *FightDiscipline) HandleEvent(event any) error {
	e, ok := event.(FightEvent)
	if !ok {
		return nil
	}

	// RESET — вне FSM
	if e.Type == EventFightReset {
		fd.Reset()
		return nil
	}

	// SETTINGS — вне FSM
	if e.Type == EventFightSettings {
		return fd.applySettings(e)
	}

	switch fd.fight.State() {

	case fight.StateIdle:
		return fd.handleIdle(e)

	case fight.StateRunning:
		return fd.handleRunning(e)

	case fight.StatePaused:
		return fd.handlePaused(e)

	case fight.StateBreak:
		return fd.handleBreak(e)

	case fight.StateFinished:
		return fight.ErrFightFinished
	}

	return nil
}

func (fd *FightDiscipline) Snapshot() map[string]any {
	red, blue := fd.scoreboard.Score()

	return map[string]any{
		"type":     "fight",
		"state":    fd.fight.State().String(),
		"round":    fd.fight.CurrentRound(),
		"rounds":   fd.fight.TotalRounds(),
		"red":      red,
		"blue":     blue,
		"warnings": fd.warnings.Snapshot(),
		"judges":   fd.scoreboard.JudgesSnapshot(), // 🔥 ВОТ ЭТО
	}
}

func (fd *FightDiscipline) Reset() {
	fd.roundTimer.Stop()
	fd.breakTimer.Stop()

	fd.fight = fight.NewFight()
	fd.scoreboard = fight.NewScoreboard()
	fd.warnings = fight.NewWarningCounter(fd.cfg.WarningsForPenalty)
}

func (fd *FightDiscipline) Realtime() <-chan RealtimeEvent {
	return fd.realtime
}

/* ================= SETTINGS ================= */

func (fd *FightDiscipline) applySettings(e FightEvent) error {
	// запрещаем менять настройки во время боя
	if fd.fight.State() == fight.StateRunning {
		return nil
	}

	if e.Rounds > 0 {
		fd.fight.SetRounds(e.Rounds)
	}

	if e.RoundDuration > 0 {
		fd.roundTimer.SetDuration(time.Duration(e.RoundDuration) * time.Second)
	}

	if e.BreakDuration > 0 {
		fd.breakTimer.SetDuration(time.Duration(e.BreakDuration) * time.Second)
	}

	return nil
}

/* ================= FSM ================= */

func (fd *FightDiscipline) handleIdle(e FightEvent) error {
	if e.Type == EventFightStart {
		fd.fight.Start()
		fd.roundTimer.Reset()
		fd.roundTimer.Start()
	}
	return nil
}

func (fd *FightDiscipline) handleRunning(e FightEvent) error {
	switch e.Type {

	case EventFightPause:
		fd.roundTimer.Stop()
		fd.breakTimer.Stop()
		fd.fight.Pause()

	case EventFightStop:
		fd.roundTimer.Stop()
		fd.breakTimer.Stop()
		fd.fight.Stop()

	case EventFightScore:
		ev := events.ScoreEvent{
			JudgeID: e.JudgeID,
			Fighter: e.Fighter,
			Points:  e.Points,
			Time:    time.Now(),
		}
		fd.scoreboard.Apply(ev)

	case EventFightWarning:
		penalty := fd.warnings.Add(e.Fighter)
		if penalty {
			fd.scoreboard.ApplyPenalty(
				e.Fighter,
				fd.cfg.PenaltyPoints,
			)
		}
	}

	return nil
}

func (fd *FightDiscipline) handlePaused(e FightEvent) error {
	if e.Type == EventFightStart {
		fd.fight.Start()
		fd.roundTimer.Reset() // 🔥 ВАЖНО
		fd.roundTimer.Start()
	}
	return nil
}

func (fd *FightDiscipline) handleBreak(e FightEvent) error {
	if e.Type == EventFightStart {
		fd.breakTimer.Stop()
		fd.fight.Start()
		fd.roundTimer.Reset()
		fd.roundTimer.Start()
	}
	return nil
}

/* ================= TIMERS ================= */

func (fd *FightDiscipline) onRoundFinished() {
	if fd.fight.CurrentRound() < fd.fight.TotalRounds() {
		fd.fight.SetState(fight.StateBreak)
		fd.breakTimer.Reset()
		fd.breakTimer.Start()
	} else {
		fd.fight.Stop()
	}
}

func (fd *FightDiscipline) onBreakFinished() {
	fd.fight.NextRound()
	fd.fight.SetState(fight.StatePaused)
}

/* ================= REALTIME ================= */

func (fd *FightDiscipline) emitTimer(mode string, rem time.Duration) {
	select {
	case fd.realtime <- RealtimeEvent{
		Type: "timer",
		Data: map[string]any{
			"mode":    mode,
			"seconds": int(rem.Seconds()),
		},
	}:
	default:
		// не блокируем таймер
	}
}
