package controllers

import (
	"flights/controllers/responses"
	"flights/errors"
	"flights/models"
	"flights/objects"
	"log"

	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

type flightCtrl struct {
	model *models.FlightsM
}

func InitFlights(r *mux.Router, model *models.FlightsM) {
	ctrl := &flightCtrl{model}
	r.HandleFunc("/flights", ctrl.getAll).Methods("GET")
	r.HandleFunc("/flights/{flightNumber}", ctrl.get).Methods("GET")
}

func (ctrl *flightCtrl) getAll(w http.ResponseWriter, r *http.Request) {
	queryParams := r.URL.Query()
	page, _ := strconv.Atoi(queryParams.Get("page"))
	page_size, _ := strconv.Atoi(queryParams.Get("size"))
	items := ctrl.model.GetAll(page, page_size)

	log.Printf("Get All flights %s\n", items)

	data := &objects.PaginationResponse{
		Page:          page,
		PageSize:      page_size,
		TotalElements: len(items),
		Items:         objects.ToFilghtResponses(items),
	}

	responses.JsonSuccess(w, data)
}

func (ctrl *flightCtrl) get(w http.ResponseWriter, r *http.Request) {
	urlParams := mux.Vars(r)
	flight_number := urlParams["flightNumber"]

	data, err := ctrl.model.Find(flight_number)
	switch err {
	case nil:
		responses.JsonSuccess(w, data.ToFilghtResponse())
	case errors.RecordNotFound:
		responses.RecordNotFound(w, flight_number)
	default:
		responses.InternalError(w)
	}
}
