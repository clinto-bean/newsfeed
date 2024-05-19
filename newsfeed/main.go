package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"

	"github.com/clinto-bean/newsfeed/internal/database"
	godotenv "github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

type apiConfig struct {
	DB *database.Queries
}

func main(){
	// Attempt to load environment variables
	log.Println("-- Attempting to load environment variables.")
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Could not load environment variables: %v", err.Error())
	}
	port := os.Getenv("PORT")
	dbUrl := os.Getenv("PSQL_DB_URL")
	log.Println("- Loaded environment variables.")

	// Create psql db connection
	log.Println("-- Attempting to connect to DB.")
	db, err := sql.Open("postgres", dbUrl)
	if err != nil {
		log.Fatalf("Could not connect to database: %v", err.Error())
	}
	dbQueries := database.New(db)
	log.Println("- Connected to DB.")

	// Create a new apiConfig
	log.Println("-- Configuring API.")
	apiCfg := apiConfig{
		DB: dbQueries,
	}
	log.Println("- API Successfully configured.")

	// create a new servemux
	log.Println("-- Creating server multiplexer.")
	mux := http.NewServeMux()
	mux.HandleFunc("GET /v1/readiness", apiCfg.handlerReadiness)
	mux.HandleFunc("GET /v1/err", apiCfg.handlerError)
	mux.HandleFunc("POST /v1/users", apiCfg.handlerCreateUsers)
	mux.HandleFunc("GET /v1/users", apiCfg.middlewareAuth(apiCfg.handlerGetUser))
	mux.HandleFunc("POST /v1/feeds", apiCfg.middlewareAuth(apiCfg.handlerCreateFeeds))

	// create cors mux for resource sharing
	log.Println("-- Creating middleware.")
	corsMux := middlewareCors(mux)

	srv := &http.Server{
		Addr: ":" + port,
		Handler: corsMux,
	}
	
	log.Println("====================")
	log.Printf("Server running on port %v!", port)
	log.Println("====================")
	log.Fatal(srv.ListenAndServe())
}

