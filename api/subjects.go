package api

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/aethiopicuschan/tsubo/subject"
)

func getSubjectsAPI(w http.ResponseWriter, r *http.Request) {
	board := r.URL.Query().Get("board")
	sj, err := subject.Get(board)
	if err != nil {
		responseError(w, r, http.StatusBadGateway, err)
		return
	}
	json, err := json.MarshalIndent(sj, "", "  ")
	if err != nil {
		responseError(w, r, http.StatusInternalServerError, err)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	fmt.Fprintf(w, string(json))
}

func SubjectsAPI(w http.ResponseWriter, r *http.Request) {
	accessLog(r)
	if r.Method == "GET" {
		getSubjectsAPI(w, r)
	} else {
		responseError(w, r, http.StatusMethodNotAllowed, nil)
	}
}
