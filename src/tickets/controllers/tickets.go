package controllers

import (
	"encoding/json"
	"log"
	"tickets/controllers/responses"
	"tickets/errors"
	"tickets/models"
	"tickets/objects"

	"net/http"

	"github.com/gorilla/mux"
)

type filghtCtrl struct {
	model *models.TicketsM
}

func InitTickets(r *mux.Router, model *models.TicketsM) {
	ctrl := &filghtCtrl{model}
	r.HandleFunc("/tickets", ctrl.fetch).Methods("GET")
	r.HandleFunc("/tickets", ctrl.create).Methods("POST")
	r.HandleFunc("/tickets/{ticketUid}", ctrl.get).Methods("GET")
	r.HandleFunc("/tickets/{ticketUid}", ctrl.delete).Methods("DELETE")
	r.HandleFunc("/manage/health", ctrl.GetHealth).Methods("GET")
}

func (h *filghtCtrl) GetHealth(w http.ResponseWriter, r *http.Request) {

	w.WriteHeader(http.StatusOK)
}

func (ctrl *filghtCtrl) fetch(w http.ResponseWriter, r *http.Request) {
	token := RetrieveToken(w, r)
	if token == nil {
		log.Printf("failed to RetrieveToken")
		return
	}
	data := ctrl.model.Fetch(token.Subject)
	responses.JsonSuccess(w, data)
}

func (ctrl *filghtCtrl) create(w http.ResponseWriter, r *http.Request) {
	req_body := new(objects.CreateRequest)
	json.NewDecoder(r.Body).Decode(req_body)
	token := RetrieveToken(w, r)
	if token == nil {
		log.Printf("failed to RetrieveToken")
		return
	}
	ticket, _ := ctrl.model.Create(token.Subject, req_body.FlightNumber, req_body.Price)
	responses.JsonSuccess(w, ticket)
}

func (ctrl *filghtCtrl) get(w http.ResponseWriter, r *http.Request) {
	urlParams := mux.Vars(r)
	ticket_uid := urlParams["ticketUid"]

	data, err := ctrl.model.Find(ticket_uid)
	switch err {
	case nil:
		responses.JsonSuccess(w, data)
	case errors.RecordNotFound:
		responses.RecordNotFound(w, ticket_uid)
	default:
		responses.InternalError(w)
	}
}

func (ctrl *filghtCtrl) delete(w http.ResponseWriter, r *http.Request) {
	urlParams := mux.Vars(r)
	ticket_uid := urlParams["ticketUid"]

	log.Printf("ticket_uid tickets", ticket_uid)

	switch ctrl.model.Delete(ticket_uid) {
	case nil:
		responses.SuccessTicketDeletion(w)
	case errors.RecordNotFound:
		responses.RecordNotFound(w, ticket_uid)
	default:
		responses.InternalError(w)
	}
}
