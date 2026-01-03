package fight

import "tkd-judge/internal/events"

type WarningCounter struct {
	red  int
	blue int
}

func NewWarningCounter() *WarningCounter {
	return &WarningCounter{}
}

func (w *WarningCounter) Add(fighter events.Fighter) (penalty bool) {
	switch fighter {
	case events.FighterRed:
		w.red++
		if w.red%3 == 0 {
			return true
		}
	case events.FighterBlue:
		w.blue++
		if w.blue%3 == 0 {
			return true
		}
	}
	return false
}

func (w *WarningCounter) Count() (red, blue int) {
	return w.red, w.blue
}
