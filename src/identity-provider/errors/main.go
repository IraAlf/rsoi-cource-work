package errors

import "errors"

var (
	RecordNotFound       = errors.New("Record was not found")
	InvalidRequest       = errors.New("Invalid rerquest")
	DatabaseWritingError = errors.New("Error while writing to DB")
	UnknownError         = errors.New("Unknown error was happened")
)
