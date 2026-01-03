package events

import "time"

type WarningEvent struct {
	Fighter Fighter
	Time    time.Time
}
