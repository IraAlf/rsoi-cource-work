package controllers

import (
	"fmt"
	"gateway/models"
	"gateway/utils"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/rs/cors"
)

func initControllers(r *mux.Router, models *models.Models) {
	r.Use(utils.LogHandler)
	api1_r := r.PathPrefix("/api/v1/").Subrouter()

	api1_r_noauth := api1_r.NewRoute().Subrouter()
	api1_r_noauth.Use(func(next http.Handler) http.Handler {
		return utils.RequestStatMiddleware(next, models.Kafka.Topic, models.Kafka.Producer)
	})
	InitAuth(api1_r_noauth, models.Client, models.Privileges)

	api1_r_auth := api1_r.NewRoute().Subrouter()
	//	api1_r_auth.Use(JwtAuthentication)
	api1_r_auth.Use(func(next http.Handler) http.Handler {
		return utils.RequestStatMiddleware(next, models.Kafka.Topic, models.Kafka.Producer)
	})

	InitFlights(api1_r_auth, models.Flights)
	InitPrivileges(api1_r_auth, models.Privileges)
	InitTickets(api1_r_auth, models.Tickets)
	InitStatistics(api1_r_auth, models.Statistics)
}

func InitRouter() *mux.Router {
	router := mux.NewRouter()
	models := models.InitModels()

	initControllers(router, models)
	return router
}

func RunRouter(r *mux.Router, port uint16) error {
	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"http://localhost:3000", "http://frontend-service:3000"},
		AllowCredentials: true,
		AllowedMethods: []string{
			http.MethodGet,
			http.MethodPost,
			http.MethodPut,
			http.MethodPatch,
			http.MethodDelete,
			http.MethodOptions,
			http.MethodHead,
		},
		AllowedHeaders: []string{"*"},
	})
	handler := c.Handler(r)
	return http.ListenAndServe(fmt.Sprintf("0.0.0.0:%d", port), handler)
}
