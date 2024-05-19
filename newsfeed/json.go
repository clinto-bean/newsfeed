package main

import (
	"encoding/json"
	"log"
	"net/http"
)


func respondWithJSON(w http.ResponseWriter, status int, payload interface{}) {
	w.Header().Set("Content-Type", "application/json")
	log.Println("apiCfg.respondWithJSON: Payload received.")
	log.Println("apiCfg.respondWithJSON: Marshalling payload.")
	dat, err := json.Marshal(payload)
	if err != nil {
		log.Printf("apiCfg.respondWithJSON: Could not marshal payload")
		w.WriteHeader(500)
		return
	}
	log.Println("apiCfg.respondWithJSON: Successfully marshalled payload.")
	w.WriteHeader(status)
	w.Write(dat)
	log.Println("apiCfg.respondWithJSON: Response sent.")
	
}

func respondWithError(w http.ResponseWriter, status int, msg string) {
	
	respondWithJSON(w, status, map[string]string{"error":msg})
	log.Printf("apiCfg.respondWithError: response sent: %v\n", msg)
	
}