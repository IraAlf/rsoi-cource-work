package controllers

import (
	"gateway/controllers/responses"
	"gateway/errors"
	"gateway/models"
	"gateway/objects"
	"strconv"

	"net/http"

	"github.com/gorilla/mux"
)

var N = 0

var cache map[string]*objects.FlightResponse

type flightCtrl struct {
	flights *models.FlightsM
}

func InitFlights(r *mux.Router, flights *models.FlightsM) {
	ctrl := &flightCtrl{flights}
	r.HandleFunc("/flights", ctrl.fetch).Methods("GET")
	r.HandleFunc("/flights/{flightNumber}", ctrl.get).Methods("GET")
}

func (ctrl *flightCtrl) fetch(w http.ResponseWriter, r *http.Request) {
	queryParams := r.URL.Query()
	page, _ := strconv.Atoi(queryParams.Get("page"))
	page_size, _ := strconv.Atoi(queryParams.Get("size"))
	data := ctrl.flights.Fetch(page, page_size, r.Header.Get("Authorization"))
	responses.JsonSuccess(w, data)
}

func (ctrl *flightCtrl) get(w http.ResponseWriter, r *http.Request) {
	urlParams := mux.Vars(r)
	flight_number := urlParams["flightNumber"]

	if N != 0 {
		inerr := ctrl.flights.ManageFlight()
		if inerr == nil {
			N = 0
		}
		if N != 0 && cache[flight_number] != nil {
			responses.JsonSuccess(w, cache[flight_number])
			return
		} else if N != 0 {
			responses.InternalError(w)
			return
		}
	}
	data, err := ctrl.flights.Find(flight_number, r.Header.Get("Authorization"))
	switch err {
	case nil:
		{
			N = 0
			responses.JsonSuccess(w, data)
			cache[flight_number] = data
		}
	case errors.FlightNotFound:
		responses.RecordNotFound(w, flight_number)
	default:
		{
			for N < 10 {
				data, err := ctrl.flights.Find(flight_number, r.Header.Get("Authorization"))
				if err != nil {
					N += 1
					continue
				} else {
					N = 0
					responses.JsonSuccess(w, data)
					break
				}
			}
			if N == 10 {
				if cache[flight_number] != nil {
					responses.JsonSuccess(w, cache[flight_number])
				}
			}
		}
	}
}
