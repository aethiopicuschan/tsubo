package api

import (
	"encoding/json"
	"errors"
	"net/http"
)

type ErrorObject struct {
	Message string `json:"message"`
}

func newErrorObject(err error) (eo ErrorObject) {
	eo.Message = err.Error()
	return
}

func responseError(w http.ResponseWriter, r *http.Request, code int, err error) {
	if err != nil {
		errorLog(err)
	} else {
		errorLog(errors.New(http.StatusText(code)))
	}
	eo := newErrorObject(errors.New(http.StatusText(code)))
	json, _ := json.MarshalIndent(eo, "", "  ")
	w.Header().Set("Content-Type", "application/json")
	http.Error(w, string(json), code)
}

func NotFound(w http.ResponseWriter, r *http.Request) {
	accessLog(r)
	responseError(w, r, http.StatusNotFound, nil)
}
