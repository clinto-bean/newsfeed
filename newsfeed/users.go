package main

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	database "github.com/clinto-bean/newsfeed/internal/database"
	uuid "github.com/google/uuid"
)

type User struct {
		ID uuid.UUID `json:"id"`
		CreatedAt time.Time `json:"created_at"`
		UpdatedAt time.Time `json:"updated_at"`
		Name string `json:"name"`
		APIKey string `json:"api_key"`
}

func (cfg *apiConfig) handlerCreateUsers(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Name string `json:"name"`
	}
	LogRequest(r)

	// decode request object

	log.Println("apiConfig.handlerCreateUsers: Attempting to decode request.")
	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		log.Printf("apiConfig.handlerCreateUsers: Could not decode request: %v\n", err.Error())
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	log.Println("apiConfig.handlerCreateUsers: Successfully decoded request parameters.")

	// define parameters for creating new user

	newUserParams := database.CreateUserParams{
        ID:        uuid.New(),
        CreatedAt: time.Now(),
        UpdatedAt: time.Now(),
        Name:      params.Name,
    }

	// attempt to create a new user

    newUser, err := cfg.DB.CreateUser(
        r.Context(),
        newUserParams,
    )

    if err != nil {
        log.Printf("apiConfig.handlerCreateUsers: Could not create new user: %v\n", err.Error())
        respondWithError(w, http.StatusInternalServerError, "Could not create new user")
        return
    }

	log.Printf("apiConfig.handlerCreateUsers: New User created: %v", newUser.Name)

	// respond to request

	respondWithJSON(w, http.StatusOK, newUser)
	
}

func (cfg *apiConfig) handlerGetUser(w http.ResponseWriter, r *http.Request, resource interface{}){

	// ensure resource passed is a User, if not, respond with error
	user, ok := resource.(database.User)
	if !ok {
		log.Println("handlerGetUser: could not assert resource as user")
        respondWithError(w, http.StatusInternalServerError, "Internal Server Error")
        return
	}

	log.Printf("apiConfig.handlerGetUser: Located user: %v\n", user.Name)

	respondWithJSON(w, http.StatusOK, User{
		ID: user.ID,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
		Name: user.Name,
		APIKey: user.Apikey,
	})
}
