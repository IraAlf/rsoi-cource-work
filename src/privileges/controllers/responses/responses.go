package responses

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type validationErrors struct{}
type validationErrorResponse struct {
	Message string           `json:"message"`
	Errors  validationErrors `json:"errors"`
}

func BadRequest(w http.ResponseWriter, msg string) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusBadRequest)
	json.NewEncoder(w).Encode(msg)
}

func ValidationErrorResponse(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusBadRequest)

	resp := &validationErrorResponse{"Request validation failed", validationErrors{}}
	json.NewEncoder(w).Encode(resp)
}

func RecordNotFound(w http.ResponseWriter, recType string) {
	msg := fmt.Sprintf("Not found %s for ID", recType)
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusNotFound)
	json.NewEncoder(w).Encode(msg)
}

func TextSuccess(w http.ResponseWriter, msg string) {
	w.Header().Set("Content-Type", "text/plain; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(msg)
}

func JsonSuccess(w http.ResponseWriter, data interface{}) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.Header().Set("Response-Code", "00")
	w.Header().Set("Response-Desc", "Success")

	json.NewEncoder(w).Encode(data)
}

func SuccessCreation(w http.ResponseWriter, location string) {
	w.Header().Set("Content-Type", "text/plain; charset=UTF-8")
	w.Header().Set("Location", location)
	w.WriteHeader(http.StatusCreated)
}

func SuccessTicketDeletion(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "text/plain; charset=UTF-8")
	w.WriteHeader(http.StatusNoContent)
	json.NewEncoder(w).Encode("Возврат билета успешно выполнен")
}

func TokenIsMissing(w http.ResponseWriter) {
	msg := "Missing auth token"
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusUnauthorized)
	json.NewEncoder(w).Encode(msg)
}

func JwtAccessDenied(w http.ResponseWriter) {
	msg := "jwt-token is not valid"
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusUnauthorized)
	json.NewEncoder(w).Encode(msg)
}

func TokenExpired(w http.ResponseWriter) {
	msg := "jwt-token expired"
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusUnauthorized)
	json.NewEncoder(w).Encode(msg)
}
