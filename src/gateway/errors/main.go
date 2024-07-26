package errors

import (
	"errors"
)

var (
	FlightNotFound  = errors.New("Flight is not found")
	ForbiddenTicket = errors.New("Forbidden ticket for this user")
)
