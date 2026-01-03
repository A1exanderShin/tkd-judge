package config

import "time"

type Config struct {
	RoundDuration      time.Duration
	AntiClick          time.Duration
	JudgesCount        int
	WarningsForPenalty int
	PenaltyPoints      int
}

func Default() Config {
	return Config{
		RoundDuration:      120 * time.Second,
		AntiClick:          300 * time.Millisecond,
		JudgesCount:        4,
		WarningsForPenalty: 3,
		PenaltyPoints:      -1,
	}
}
