package controllers

import (
	"fmt"
	"net/http"
	"statistics/models"
	"statistics/objects"
	"statistics/utils"
	"time"

	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
	"github.com/rs/cors"
)

type requestCtrl struct {
	statistics *models.StatisticsM
}

func initControllers(r *mux.Router, models *models.Models) {
	r.Use(utils.LogHandler)

	api1_r := r.PathPrefix("/api/v1/").Subrouter()
	api1_r.Use(JwtAuthentication)

	ctrl := &requestCtrl{models.Statistics}
	api1_r.HandleFunc("/requests", ctrl.fetch).Methods("GET")
}

func InitRouter(db *gorm.DB) *mux.Router {
	router := mux.NewRouter()
	models := models.InitModels(db)

	initControllers(router, models)
	return router
}

func RunRouter(r *mux.Router, port uint16) error {
	c := cors.New(cors.Options{})
	handler := c.Handler(r)
	return http.ListenAndServe(fmt.Sprintf("0.0.0.0:%d", port), handler)
}

func (ctrl *requestCtrl) fetch(w http.ResponseWriter, r *http.Request) {
	queryParams := r.URL.Query()
	begin_time, err := time.Parse(time.RFC3339, queryParams.Get("begin_time"))
	if err != nil {
		http.Error(w, "bad begin_time format", http.StatusBadRequest)
		return
	}
	end_time, err := time.Parse(time.RFC3339, queryParams.Get("end_time"))
	if err != nil {
		http.Error(w, "bad end_time format", http.StatusBadRequest)
		return
	}

	requests, err := ctrl.statistics.Fetch(begin_time, end_time)
	if err != nil {
		http.Error(w, "failed to fetch statistics", http.StatusInternalServerError)
		return
	}
	JsonSuccess(w, &objects.FetchResponse{Reqests: requests})
}
