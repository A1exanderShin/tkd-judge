package fight

import "tkd-judge/internal/events"

type Scoreboard struct {
	red  int
	blue int
}

func NewScoreboard() *Scoreboard {
	return &Scoreboard{}
}

func (s *Scoreboard) Apply(event events.ScoreEvent) {
	switch event.Fighter {
	case events.FighterRed:
		s.red += event.Points
	case events.FighterBlue:
		s.blue += event.Points
	}
}

func (s *Scoreboard) Score() (red, blue int) {
	return s.red, s.blue
}
