package fight

import "tkd-judge/internal/events"

type JudgeScore struct {
	Red  int
	Blue int
}

type Scoreboard struct {
	red    int
	blue   int
	judges map[int]*JudgeScore
}

func NewScoreboard() *Scoreboard {
	return &Scoreboard{
		judges: make(map[int]*JudgeScore),
	}
}

// принять конкретное доменное событие
func (s *Scoreboard) Apply(event events.ScoreEvent) {
	switch event.Fighter {
	case events.FighterRed:
		s.red += event.Points
	case events.FighterBlue:
		s.blue += event.Points
	}

	js, ok := s.judges[event.JudgeID]
	if !ok {
		js = &JudgeScore{}
		s.judges[event.JudgeID] = js
	}

	switch event.Fighter {
	case events.FighterRed:
		js.Red += event.Points
	case events.FighterBlue:
		js.Blue += event.Points
	}
}

func (s *Scoreboard) Score() (red, blue int) {
	return s.red, s.blue
}

func (s *Scoreboard) Reset() {
	s.red = 0
	s.blue = 0
}

func (s *Scoreboard) JudgesSnapshot() []map[string]int {
	out := make([]map[string]int, 0, len(s.judges))

	for id, score := range s.judges {
		out = append(out, map[string]int{
			"id":   id,
			"red":  score.Red,
			"blue": score.Blue,
		})
	}

	return out
}

func (s *Scoreboard) ApplyPenalty(fighter events.Fighter, points int) {
	switch fighter {
	case events.FighterRed:
		s.red += points
	case events.FighterBlue:
		s.blue += points
	}
}
