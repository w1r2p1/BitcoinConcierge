package utils

import (
	"encoding/json"
	"net/http"
)

type appError struct {
	Error      string `json:"error"`
	Message    string `json:"message"`
	HttpStatus int    `json:"status"`
}

func DisplayError(w http.ResponseWriter, handlerError error, message string, code int) {
	errObj := appError{
		Error:      handlerError.Error(),
		Message:    message,
		HttpStatus: code,
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(code)
	if j, err := json.Marshal(errObj); err != nil {
		w.Write(j)
	}
}
