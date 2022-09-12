package api

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/aethiopicuschan/tsubo/bbsmenu"
)

func getBBSMenuAPI(w http.ResponseWriter, r *http.Request) {
	bm, err := bbsmenu.Get("https://menu.5ch.net")
	if err != nil {
		responseError(w, r, http.StatusBadGateway, err)
		return
	}
	json, err := json.MarshalIndent(bm, "", "  ")
	if err != nil {
		responseError(w, r, http.StatusInternalServerError, err)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	fmt.Fprintf(w, string(json))
}

func BBSMenuAPI(w http.ResponseWriter, r *http.Request) {
	accessLog(r)
	if r.Method == "GET" {
		getBBSMenuAPI(w, r)
	} else {
		responseError(w, r, http.StatusMethodNotAllowed, nil)
	}
}
