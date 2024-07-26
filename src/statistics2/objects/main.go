package objects

import (
	_ "encoding/json"
	"time"
)

type RequestStat struct {
	Id           int           `gorm:"primary_key;index"`
	Path         string        `json:"path" gorm:"not null"`
	ResponceCode int           `json:"responceCode" gorm:"not null"`
	Method       string        `json:"method" gorm:"not null"`
	StartedAt    time.Time     `json:"startedAt" gorm:"not null"`
	FinishedAt   time.Time     `json:"finishedAt" gorm:"not null"`
	Duration     time.Duration `json:"duration" gorm:"not null"`
	UserName     string        `json:"userName,omitempty"`
}

func (RequestStat) TableName() string {
	return "statistics"
}

type FetchResponse struct {
	Reqests []RequestStat `json:"requests"`
}
