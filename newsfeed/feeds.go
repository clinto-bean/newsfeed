package main

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/clinto-bean/newsfeed/internal/database"
	uuid "github.com/google/uuid"
)

type Feed struct {
		ID uuid.UUID `json:"id"`
		CreatedAt time.Time `json:"created_at"`
		UpdatedAt time.Time `json:"updated_at"`
		Name string `json:"name"`
		Url string `json:"url"`
		UserID uuid.UUID `json:"user_id"`
}

func (cfg *apiConfig) handlerCreateFeeds(w http.ResponseWriter, r *http.Request, user interface{}) {
	type parameters struct {
		Name string `json:"name"`
		Url string `json:"url"`
	}

	dbUser, ok := user.(database.User)
	if !ok {
		log.Println("apiConfig.handlerCreateFeeds: user was not authorized")
		respondWithError(w, http.StatusUnauthorized, "apiConfig.handlerCreateFeeds: Please sign in")
		return
	}

	log.Printf("%s: Request received: %v %v", getFunctionName(cfg), r.Method, r.URL.String())

	// decode request object

	log.Println("apiConfig.handlerCreateFeeds: Attempting to decode request.")
	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		log.Printf("apiConfig.handlerCreateFeeds: Could not decode request: %v\n", err.Error())
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	log.Println("apiConfig.handlerCreateFeeds: Successfully decoded request parameters.")

	newFeedParams := database.CreateFeedParams{
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Name: params.Name,
		Url: params.Url,
		UserID: dbUser.ID,
	}
	log.Println("apiConfig.handlerCreateFeeds: Creating feed in database")
	
	createdFeed, err := cfg.DB.CreateFeed(r.Context(), newFeedParams)
	if err != nil {
		log.Println("apiCfg.handlerCreateFeeds: could not create user in DB")
		respondWithError(w, http.StatusInternalServerError, "Couldn't save feed")
		return
	}
	log.Println("apiConfig.handlerCreateFeeds: Feed properly saved to database")

	respondWithJSON(w, http.StatusCreated, createdFeed)

}