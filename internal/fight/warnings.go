package fight

import "tkd-judge/internal/events"

type WarningCounter struct {
	red   int
	blue  int
	limit int
}

func NewWarningCounter(limit int) *WarningCounter {
	return &WarningCounter{limit: limit}
}

// добавляет штрафы
func (w *WarningCounter) Add(fighter events.Fighter) (penalty bool) {
	switch fighter {
	// проверка кратности, каждое 3-е предупреждение = штраф
	case events.FighterRed:
		w.red++
		if w.red%w.limit == 0 {
			return true
		}
	case events.FighterBlue:
		w.blue++
		if w.blue%w.limit == 0 {
			return true
		}
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
