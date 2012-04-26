package controller

import (
	"net/http"
)

type Response struct {
	Status int
	Header http.Header
	Message []byte
}

func (this Response) ToRack() (int,http.Header,[]byte) {
	return this.Status,this.Header,this.Message
}

func FromRack(status int, header http.Header, message []byte) (result Response) {
	result.Status = status
	result.Header = header
	result.Message = message
	return
}

func BlankResponse(status int) Response {
	return FromRack(status,make(http.Header),[]byte(""))
}

func NotFound() Response{
	return BlankResponse(http.StatusNotFound)
}

func NotAuthorized() Response{
	return BlankResponse(http.StatusUnauthorized)
}