package repository

import (
	"flights/errors"
	"flights/objects"

	"github.com/jinzhu/gorm"
)

type FlightsRep interface {
	GetAll(page int, page_size int) []objects.Flight
	Find(flight_number string) (*objects.Flight, error)
}

type PGFlightsRep struct {
	db *gorm.DB
}

func NewPGFlightsRep(db *gorm.DB) *PGFlightsRep {
	return &PGFlightsRep{db}
}

func paginate(page int, page_size int) func(*gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		offet := (page - 1) * page_size
		return db.Offset(offet).Limit(page_size)
	}
}

func (rep *PGFlightsRep) GetAll(page int, page_size int) []objects.Flight {
	temp := []objects.Flight{}
	rep.db.
		Scopes(paginate(page, page_size)).
		Model(&objects.Flight{}).
		Preload("FromAirport").
		Preload("ToAirport").
		Find(&temp)

	return temp
}

func (rep *PGFlightsRep) Find(flight_number string) (*objects.Flight, error) {
	temp := new(objects.Flight)
	err := rep.db.
		Where(&objects.Flight{FlightNumber: flight_number}).
		Preload("FromAirport").
		Preload("ToAirport").
		First(temp).
		Error
	switch err {
	case nil:
		break
	case gorm.ErrRecordNotFound:
		temp, err = nil, errors.RecordNotFound
	default:
		temp, err = nil, errors.UnknownError
	}

	return temp, err
}
