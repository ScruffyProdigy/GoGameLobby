/*
Error Handler Package recovers from any downstream errors, and will return a 500 status, and set the body the Error Message
*/
package errorhandler

import (
	"../rack"
	"net/http"
)

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

func ErrorHandler(r *http.Request, vars rack.Vars, next rack.NextFunc) (status int, header http.Header, message []byte) {
	defer func() {
		rec := recover()
		if rec != nil {
			status = http.StatusInternalServerError
			message = []byte(getErrorString(rec))
			if header == nil {
				header = make(http.Header)
			}
		}
	}()
	return next()
}
