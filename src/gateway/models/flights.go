package models

import (
	"context"
	"encoding/json"
	"fmt"
	"gateway/errors"
	"gateway/objects"
	"gateway/utils"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

type FlightsM struct {
	client *http.Client
}

func NewFlightsM(client *http.Client) *FlightsM {
	return &FlightsM{client: client}
}

func (model *FlightsM) Fetch(page int, page_size int, authHeader string) *objects.PaginationResponse {
	req, _ := http.NewRequest("GET", fmt.Sprintf("%s/api/v1/flights", utils.Config.Endpoints.Flights), nil)
	q := req.URL.Query()
	q.Add("page", fmt.Sprintf("%d", page))
	q.Add("size", fmt.Sprintf("%d", page_size))
	req.URL.RawQuery = q.Encode()
	req.Header.Add("Authorization", authHeader)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	req = req.WithContext(ctx)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		panic("client: error making http request\n")
	}

	data := &objects.PaginationResponse{}
	body, _ := ioutil.ReadAll(resp.Body)
	json.Unmarshal(body, data)

	log.Printf("flights: %v", data)
	return data
}

func (model *FlightsM) Find(flight_number string, authHeader string) (*objects.FlightResponse, error) {
	req, _ := http.NewRequest("GET", fmt.Sprintf("%s/api/v1/flights/%s", utils.Config.Endpoints.Flights, flight_number), nil)
	req.Header.Add("Authorization", authHeader)
	resp, err := model.client.Do(req)
	if err != nil {
		return nil, err
	} else if resp.StatusCode == http.StatusNotFound {
		return nil, errors.FlightNotFound
	} else {
		data := &objects.FlightResponse{}
		body, _ := ioutil.ReadAll(resp.Body)
		json.Unmarshal(body, data)
		return data, nil
	}
}

func (model *FlightsM) ManageFlight() error {

	req, _ := http.NewRequest("GET", fmt.Sprintf("%s/api/v1/flights/manage/health", utils.Config.Endpoints.Flights), nil)
	_, err := model.client.Do(req)
	if err != nil {
		return err
	} else {
		return nil
	}
}
