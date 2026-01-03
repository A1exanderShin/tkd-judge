package protocol

import (
	"time"

	"tkd-judge/internal/events"
)

type Protocol struct {
	GeneratedAt time.Time `json:"generated_at"`

	State string `json:"state"`

	Score struct {
		Red  int `json:"red"`
		Blue int `json:"blue"`
	} `json:"score"`

	Warnings struct {
		Red  int `json:"red"`
		Blue int `json:"blue"`
	} `json:"warnings"`

	Events []any `json:"events"`
}

type ScoreLog struct {
	Type    string         `json:"type"`
	JudgeID int            `json:"judge_id"`
	Fighter events.Fighter `json:"fighter"`
	Points  int            `json:"points"`
	Time    time.Time      `json:"time"`
}

type WarningLog struct {
	Type    string         `json:"type"`
	Fighter events.Fighter `json:"fighter"`
	Time    time.Time      `json:"time"`
}
