package main

import (
	"log"
	"net/http"
)

type authHandler func(
	http.ResponseWriter, 
	*http.Request, 
	interface{},
)

func middlewareCors(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html;encoding=utf-8")
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, DELETE, PUT, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "*")
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func (cfg *apiConfig) middlewareAuth(handler authHandler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		// parse api key

		authorization := r.Header.Get("Authorization")
		log.Printf("apikey: %v\n", authorization)
		if len(authorization) < 6 || authorization[:6] != "apikey" {
			log.Println("apiCfg.middlewareAuth: bad authorization request")
			respondWithError(w, http.StatusBadRequest, "Please use your API key")
			return
		}

		apiKey := authorization[7:]

		// extract apikey and attempt to locate resource

		var resource interface{}
        var err error

		switch r.URL.Path {
		case "/v1/users":
			resource, err = cfg.DB.GetUser(r.Context(), apiKey)
			if err != nil {
				log.Println("apiCfg.middlewareAuth: could not locate user given parameters")
				respondWithError(w, http.StatusNotFound, "Could not locate user")
				return
			}
		default:
			log.Println("apiCfg.middlewareAuth: unknown path")
			respondWithError(w, http.StatusNotFound, "Resource not found")
			return
		}

		// invoke handler with writer, request and user

		handler(w, r, resource)
	}
}

