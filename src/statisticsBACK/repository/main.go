package repository

import (
	"log"
	"statistics/objects"
	"time"

	"github.com/jinzhu/gorm"
)

type StatisticsRep interface {
	Create(*objects.RequestStat) error
	Fetch(beginTime time.Time, endTime time.Time) ([]objects.RequestStat, error)
}

type PGStatisticsRep struct {
	db *gorm.DB
}

func NewPGStatisticsRep(db *gorm.DB) *PGStatisticsRep {
	return &PGStatisticsRep{db}
}

func (rep *PGStatisticsRep) Create(statistics *objects.RequestStat) error {
	log.Println("writing statistics to db")
	return rep.db.Create(statistics).Error
}

func (rep *PGStatisticsRep) Fetch(beginTime time.Time, endTime time.Time) ([]objects.RequestStat, error) {
	log.Println("fetching requests from db")
	var requests []objects.RequestStat
	err := rep.db.
		Model(&objects.RequestStat{}).
		Where("statistics.started_at >= ? AND statistics.started_at <= ?", beginTime, endTime).
		Find(&requests).Error
	return requests, err
}
