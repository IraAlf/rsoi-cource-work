package models

import (
	"statistics/objects"
	"statistics/repository"
	"time"

	"github.com/jinzhu/gorm"
)

type StatisticsM struct {
	rep repository.StatisticsRep
}

type Models struct {
	Statistics *StatisticsM
}

func InitModels(db *gorm.DB) *Models {
	models := new(Models)
	models.Statistics = &StatisticsM{repository.NewPGStatisticsRep(db)}
	return models
}

func (model *StatisticsM) Fetch(beginTime time.Time, endTime time.Time) ([]objects.RequestStat, error) {
	return model.rep.Fetch(beginTime, endTime)
}
