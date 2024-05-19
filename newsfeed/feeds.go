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
		ID string `json:"id"`
		CreatedAt time.Time `json:"created_at"`
		UpdatedAt time.Time `json:"updated_at"`
		Name string `json:"name"`
		Url string `json:"url"`
		UserID string `json:"user_id"`
}

func (cfg *apiConfig) handlerCreateFeeds(w http.ResponseWriter, r *http.Request, user interface{}) {

	type parameters struct {
		Name string `json:"name"`
		Url string `json:"url"`
	}

	dbUser, ok := user.(database.User)
	if !ok {
		log.Println("apiCfg.handlerCreateFeeds: user was not authorized")
		respondWithError(w, http.StatusUnauthorized, "apiCfg.handlerCreateFeeds: Please sign in")
		return
	}

	// decode request object

	log.Println("apiCfg.handlerCreateFeeds: Attempting to decode request.")
	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		log.Printf("apiCfg.handlerCreateFeeds: Could not decode request: %v\n", err.Error())
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	log.Println("apiCfg.handlerCreateFeeds: Successfully decoded request parameters.")

	newFeedParams := database.CreateFeedParams{
		ID: uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Name: params.Name,
		Url: params.Url,
		UserID: dbUser.ID,
	}

	log.Println("apiCfg.handlerCreateFeeds: Creating feed in database")
	createdFeed, err := cfg.DB.CreateFeed(r.Context(), newFeedParams)
	log.Println()
	log.Println(cfg.DB.CreateFeed(r.Context(), newFeedParams))
	log.Println()
	if err != nil {
		log.Println("apiCfg.handlerCreateFeeds: could not create feed in DB")
		respondWithError(w, http.StatusInternalServerError, "Couldn't save feed")
		return
	}
	log.Println("apiCfg.handlerCreateFeeds: Feed properly saved to database")

	respondWithJSON(w, http.StatusCreated, createdFeed)

}