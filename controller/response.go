package controller

import (
	"net/http"
)

// to simplify the controller, instead of requiring you to type in 3 output variables for all of your control functions
// of which you will rarely interact with any of them
// we have encapsulated them in the Response struct
// each of the members is accesible if you need to access them, but we doubt you will need to do so
type Response struct {
	Status  int
	Header  http.Header
	Message []byte
}

// this function converts a controller response back into a rack response
func (this Response) ToRack() (int, http.Header, []byte) {
	return this.Status, this.Header, this.Message
}

// this function takes a rack response, and converts it into a controller response
func FromRack(status int, header http.Header, message []byte) (result Response) {
	result.Status = status
	result.Header = header
	result.Message = message
	return
}

// sometimes you just want to send back a blank response
// this will give you that
// specify the http status, and you're set to go
func BlankResponse(status int) Response {
	return FromRack(status, make(http.Header), []byte(""))
}

// specific implementation of BlankResponse for a 400 error
func BadRequest() Response {
	return BlankResponse(http.StatusBadRequest)
}

// specific implementation of BlankResponse for a 401 error
func NotAuthorized() Response {
	return BlankResponse(http.StatusUnauthorized)
}

// specific implementation of BlankResponse for a 403 error
func InsufficientAuthorization() Response {
	return BlankResponse(http.StatusUnauthorized)
}

// specific implementation of BlankResponse for a 404 error
func NotFound() Response {
	return BlankResponse(http.StatusNotFound)
}
