package controllers

import (
	"fmt"
	"gateway/controllers/responses"
	"gateway/models"
	"log"
	"time"

	"net/http"

	"github.com/gorilla/mux"
)

type statisticsCtrl struct {
	statistics *models.StatisticsM
}

func InitStatistics(r *mux.Router, statistics *models.StatisticsM) {
	ctrl := &statisticsCtrl{statistics}
	r.HandleFunc("/requests", ctrl.fetch).Methods("GET")
}

func (ctrl *statisticsCtrl) fetch(w http.ResponseWriter, r *http.Request) {
	token, err := RetrieveToken(w, r)
	if err != nil {
		log.Printf("failed to RetrieveToken: %s", err.Error())
		return
	}

	if token.Role != "admin" {
		responses.ForbiddenMsg(w, fmt.Sprintf("not allowed for %s role", token.Role))
		return
	}

	queryParams := r.URL.Query()
	log.Println(queryParams.Get("begin_time"))
	log.Println(queryParams.Get("end_time"))
	begin_time, err := time.Parse(time.RFC3339, queryParams.Get("begin_time"))
	if err != nil {
		http.Error(w, fmt.Sprintf("bad begin_time format: %s", err.Error()), http.StatusBadRequest)
		return
	}
	end_time, err := time.Parse(time.RFC3339, queryParams.Get("end_time"))
	if err != nil {
		http.Error(w, "bad end_time format", http.StatusBadRequest)
		return
	}

	data := ctrl.statistics.Fetch(begin_time, end_time, r.Header.Get("Authorization"))
	responses.JsonSuccess(w, data)
}
