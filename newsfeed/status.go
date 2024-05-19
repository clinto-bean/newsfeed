package main

import (
	"encoding/json"
	"log"
	"net/http"
)

// handlerReadiness responds to GET requests at /v1/readiness and returns either "ok" and a 200 status, or an error message and a 500 status

func (cfg *apiConfig) handlerReadiness(w http.ResponseWriter, r *http.Request) {
	type returnParams struct {
		Status string `json:"status"`
	}
	
	params := &returnParams{
		Status: "ok",
	}
	
	dat, err := json.Marshal(params)
	
	if err != nil {
		log.Printf("handlerReadiness: Could not marshal return parameters for /v1/readiness")
		return
	}

	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write(dat)
}

// handlerError responds to GET requests at /v1/err and will respond with a 500 status and the text "internal server error"

func (cfg *apiConfig) handlerError(w http.ResponseWriter, r *http.Request){
	
	type errParams struct {
		Error string `json:"error"`
	}

	params := &errParams{
		Error: "Internal Server Error",
	}

	dat, err := json.Marshal(params)

	if err != nil {
		log.Printf("handlerError: Could not marshal return parameters for /v1/readiness")
		return
	}

	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusInternalServerError)
	w.Write(dat)
}