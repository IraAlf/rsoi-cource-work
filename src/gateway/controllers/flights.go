package controllers

import (
	"gateway/controllers/responses"
	"gateway/errors"
	"gateway/models"
	"strconv"

	"net/http"

	"github.com/gorilla/mux"
)

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

	data, err := ctrl.flights.Find(flight_number, r.Header.Get("Authorization"))
	switch err {
	case nil:
		responses.JsonSuccess(w, data)
	case errors.FlightNotFound:
		responses.RecordNotFound(w, flight_number)
	default:
		responses.InternalError(w)
	}
}
