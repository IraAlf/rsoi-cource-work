package objects

import (
	_ "encoding/json"
	"time"
)

type RequestStat struct {
	Path         string        `json:"path"`
	ResponceCode int           `json:"responceCode"`
	Method       string        `json:"method"`
	StartedAt    time.Time     `json:"startedAt"`
	FinishedAt   time.Time     `json:"finishedAt"`
	Duration     time.Duration `json:"duration"`
	UserName     string        `json:"userName,omitempty"`
}

type FetchResponse struct {
	Reqests []RequestStat `json:"requests"`
}
