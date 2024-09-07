package controllers

import (
	"gateway/controllers/responses"
	"gateway/errors"
	"gateway/models"
	"gateway/objects"
	"log"

	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
)

type ticketsCtrl struct {
	tickets *models.TicketsM
}

func InitTickets(r *mux.Router, tickets *models.TicketsM) {
	ctrl := &ticketsCtrl{tickets: tickets}
	r.HandleFunc("/me", ctrl.me).Methods("GET")
	r.HandleFunc("/tickets", ctrl.fetch).Methods("GET")
	r.HandleFunc("/tickets", ctrl.post).Methods("POST")
	r.HandleFunc("/tickets/{ticketUid}", ctrl.get).Methods("GET")
	r.HandleFunc("/tickets/{ticketUid}", ctrl.delete).Methods("DELETE")
}

func (ctrl *ticketsCtrl) me(w http.ResponseWriter, r *http.Request) {
	data, err := ctrl.tickets.FetchUser(r.Header.Get("Authorization"))
	if err != nil {
		responses.InternalError(w)
	} else {
		responses.JsonSuccess(w, data)
	}
}

func (ctrl *ticketsCtrl) fetch(w http.ResponseWriter, r *http.Request) {
	data, err := ctrl.tickets.Fetch(r.Header.Get("Authorization"))
	if err != nil {
		responses.InternalError(w)
	} else {
		responses.JsonSuccess(w, data)
	}
}

func (ctrl *ticketsCtrl) post(w http.ResponseWriter, r *http.Request) {
	req_body := new(objects.TicketPurchaseRequest)
	err := json.NewDecoder(r.Body).Decode(req_body)
	if err != nil {
		responses.ValidationErrorResponse(w, err.Error())
		return
	}

	data, err := ctrl.tickets.Create(req_body.FlightNumber, r.Header.Get("Authorization"), req_body.Price, req_body.PaidFromBalance)
	if err != nil {
		responses.InternalError(w)
	} else {
		responses.JsonSuccess(w, data)
	}
}

func (ctrl *ticketsCtrl) get(w http.ResponseWriter, r *http.Request) {
	token, err := RetrieveToken(w, r)
	if err != nil {
		log.Printf("failed to RetrieveToken: %s", err.Error())
		return
	}

	urlParams := mux.Vars(r)
	ticket_uid := urlParams["ticketUid"]
	authHeader := r.Header.Get("Authorization")
	username := token.Subject
	log.Printf("username: ", username)

	data, err := ctrl.tickets.Find(ticket_uid, username, authHeader)
	switch err {
	case nil:
		responses.JsonSuccess(w, data)
	case errors.ForbiddenTicket:
		responses.Forbidden(w)
	default:
		responses.RecordNotFound(w, ticket_uid)
	}
}

func (ctrl *ticketsCtrl) delete(w http.ResponseWriter, r *http.Request) {
	token, err := RetrieveToken(w, r)
	if err != nil {
		log.Printf("failed to RetrieveToken: %s", err.Error())
		return
	}

	urlParams := mux.Vars(r)
	ticket_uid := urlParams["ticketUid"]
	authHeader := r.Header.Get("Authorization")
	username := token.Subject
	log.Printf("username: ", username)

	log.Printf("ticket_uid: ", ticket_uid)
	log.Printf("authHeader: ", authHeader)
	log.Printf("username: ", username)

	err = ctrl.tickets.Delete(ticket_uid, username, authHeader)
	switch err {
	case nil:
		responses.SuccessTicketDeletion(w)
	case errors.ForbiddenTicket:
		responses.Forbidden(w)
	default:
		responses.RecordNotFound(w, ticket_uid)
	}
}
