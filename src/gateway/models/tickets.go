package models

import (
	"bytes"
	"encoding/json"
	"fmt"
	"gateway/errors"
	"gateway/objects"
	"gateway/utils"
	"io/ioutil"
	"net/http"
)

type TicketsM struct {
	client *http.Client

	flights    *FlightsM
	privileges *PrivilegesM
}

func NewTicketsM(client *http.Client, flights *FlightsM) *TicketsM {
	return &TicketsM{
		client:  client,
		flights: flights,
	}
}

func (model *TicketsM) FetchUser(authHeader string) (*objects.UserInfoResponse, error) {
	data := new(objects.UserInfoResponse)
	tickets, err := model.fetch(authHeader)
	if err != nil {
		return nil, err
	}
	flights := model.flights.Fetch(1, 100, authHeader).Items
	data.Tickets = objects.MakeTicketResponseArr(tickets, flights)

	privilege := model.privileges.Fetch(authHeader)
	data.Privilege = objects.PrivilegeShortInfo{
		Balance: privilege.Balance,
		Status:  privilege.Status,
	}
	return data, nil
}

func (model *TicketsM) fetch(authHeader string) (objects.TicketArr, error) {
	req, _ := http.NewRequest("GET", fmt.Sprintf("%s/api/v1/tickets", utils.Config.Endpoints.Tickets), nil)
	req.Header.Set("Authorization", authHeader)
	resp, err := model.client.Do(req)
	if err != nil {
		return nil, err
	}

	data := new(objects.TicketArr)
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	json.Unmarshal(body, data)
	return *data, nil
}

func (model *TicketsM) Fetch(authHeader string) ([]objects.TicketResponse, error) {
	tickets, err := model.fetch(authHeader)
	if err != nil {
		return nil, err
	}

	flights := model.flights.Fetch(1, 100, authHeader).Items
	return objects.MakeTicketResponseArr(tickets, flights), nil
}

func (model *TicketsM) create(flight_number string, price int, authHeader string) (*objects.TicketCreateResponse, error) {
	req_body, _ := json.Marshal(&objects.TicketCreateRequest{FlightNumber: flight_number, Price: price})
	req, _ := http.NewRequest("POST", fmt.Sprintf("%s/api/v1/tickets", utils.Config.Endpoints.Tickets), bytes.NewBuffer(req_body))
	req.Header.Add("Authorization", authHeader)

	if resp, err := model.client.Do(req); err != nil {
		return nil, err
	} else {
		data := &objects.TicketCreateResponse{}
		body, _ := ioutil.ReadAll(resp.Body)
		json.Unmarshal(body, data)
		return data, nil
	}
}

func (model *TicketsM) Create(flight_number string, authHeader string, price int, from_balance bool) (*objects.TicketPurchaseResponse, error) {
	flight, err := model.flights.Find(flight_number, authHeader)
	if err != nil {
		utils.Logger.Println(err.Error())
		return nil, err
	}

	ticket, err := model.create(flight_number, price, authHeader)
	if err != nil {
		utils.Logger.Println(err.Error())
		return nil, err
	}

	privilege, err := model.privileges.AddTicket(authHeader, &objects.AddHistoryRequest{
		TicketUID:       ticket.TicketUid,
		Price:           flight.Price,
		PaidFromBalance: from_balance,
	})
	if err != nil {
		utils.Logger.Println(err.Error())
		return nil, err
	}

	return objects.NewTicketPurchaseResponse(flight, ticket, privilege), nil
}

func (model *TicketsM) find(ticket_uid string, authHeader string) (*objects.Ticket, error) {
	req, _ := http.NewRequest("GET", fmt.Sprintf("%s/api/v1/tickets/%s", utils.Config.Endpoints.Tickets, ticket_uid), nil)
	req.Header.Add("Authorization", authHeader)
	resp, err := model.client.Do(req)
	if err != nil {
		return nil, err
	} else {
		data := &objects.Ticket{}
		body, _ := ioutil.ReadAll(resp.Body)
		json.Unmarshal(body, data)
		return data, nil
	}
}

func (model *TicketsM) Find(ticket_uid string, username string, authHeader string) (*objects.TicketResponse, error) {
	ticket, err := model.find(ticket_uid, authHeader)
	if err != nil {
		return nil, err
	} else if username != ticket.Username {
		return nil, errors.ForbiddenTicket
	}

	flight, err := model.flights.Find(ticket.FlightNumber, authHeader)
	if err != nil {
		return nil, err
	} else {
		return objects.ToTicketResponce(ticket, flight), nil
	}
}

func (model *TicketsM) delete(ticket_uid string, authHeader string) error {
	req, _ := http.NewRequest("DELETE", fmt.Sprintf("%s/api/v1/tickets/%s", utils.Config.Endpoints.Tickets, ticket_uid), nil)
	req.Header.Add("Authorization", authHeader)
	_, err := model.client.Do(req)
	return err
}

func (model *TicketsM) Delete(ticket_uid string, username string, authHeader string) error {
	ticket, err := model.find(ticket_uid, authHeader)
	if err != nil {
		return err
	} else if username != ticket.Username {
		return errors.ForbiddenTicket
	}

	if err = model.delete(ticket_uid, authHeader); err != nil {
		return err
	}

	return model.privileges.DeleteTicket(authHeader, ticket_uid)
}
