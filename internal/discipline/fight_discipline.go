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
}

func NewFightDiscipline() *FightDiscipline {
	cfg := config.Default()

	fd := &FightDiscipline{
		cfg:        cfg,
		fight:      fight.NewFight(),
		scoreboard: fight.NewScoreboard(),
		warnings:   fight.NewWarningCounter(cfg.WarningsForPenalty),
	}

	fd.roundTimer = fight.NewTimer(cfg.RoundDuration)
	fd.breakTimer = fight.NewTimer(30 * time.Second)

	fd.roundTimer.OnFinished(fd.onRoundFinished)
	fd.breakTimer.OnFinished(fd.onBreakFinished)

	return fd
}

func (fd *FightDiscipline) HandleEvent(event any) error {
	e, ok := event.(FightEvent)
	if !ok {
		return nil
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
	}

	return nil
}

func (fd *FightDiscipline) handlePaused(e FightEvent) error {
	if e.Type == EventFightStart {
		fd.fight.Start()
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
	}
}

func (fd *FightDiscipline) Reset() {
	fd.roundTimer.Stop()
	fd.breakTimer.Stop()

	fd.fight = fight.NewFight()
	fd.scoreboard = fight.NewScoreboard()
	fd.warnings = fight.NewWarningCounter(fd.cfg.WarningsForPenalty)
}
