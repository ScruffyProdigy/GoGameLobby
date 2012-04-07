package routes

import (
	"../log"
	"fmt"
	"net/http"
)

type ErrorHandler func(w http.ResponseWriter, err string)

var errorhandler ErrorHandler

func SetErrorHandler(handler ErrorHandler) {
	errorhandler = handler
}

func defaultErrorHandler(w http.ResponseWriter, error_string string) {
	log.Info("500")
	w.WriteHeader(500)

	fmt.Fprint(w, "<html><head><title>Error</title></head><body><h1>Error</h1><p>", error_string, "</p></body></html>")
}

func getErrorString(rec interface{}) string {
	err, isError := rec.(error)
	str, isString := rec.(string)

	if isError {
		return err.Error()
	} else if isString {
		return str
	}
	return "Unknown Error"
}

func handleErrors(w http.ResponseWriter) {
	if rec := recover(); rec != nil {
		errorhandler(w, getErrorString(rec))
	}
}

func init() {
	SetErrorHandler(defaultErrorHandler)
}
