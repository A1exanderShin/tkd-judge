package fight

import "tkd-judge/internal/events"

type WarningCounter struct {
	red   int
	blue  int
	limit int
}

func NewWarningCounter(limit int) *WarningCounter {
	if limit <= 0 {
		limit = 1
	}
	return &WarningCounter{limit: limit}
}

// Add добавляет предупреждение.
// Возвращает true, если нужно применить штраф.
func (w *WarningCounter) Add(fighter events.Fighter) (penalty bool) {
	switch fighter {
	case events.FighterRed:
		w.red++
		return w.red%w.limit == 0

	case events.FighterBlue:
		w.blue++
		return w.blue%w.limit == 0
	}
	return false
}

func (w *WarningCounter) Count() (red, blue int) {
	return w.red, w.blue
}

func (w *WarningCounter) Snapshot() map[string]int {
	return map[string]int{
		"red":  w.red,
		"blue": w.blue,
	}
}

func (w *WarningCounter) Reset() {
	w.red = 0
	w.blue = 0
}
