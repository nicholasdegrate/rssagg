package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"

	"github.com/blindbat/rssagg/internal/database"
	"github.com/go-chi/chi"
	"github.com/go-chi/cors"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

type apiConfig struct {
	DB *database.Queries
}

func main() {
	godotenv.Load()
	
	db_url := os.Getenv("DB_URL")
	if db_url == "" {
		log.Fatal("DB_URL is not found in environment")
	}
	
	port := os.Getenv("PORT")
	if port == "" {
		log.Fatal("PORT is not found in environment")
	}
	conn, err := sql.Open("postgres", db_url)

	if err != nil {
		log.Fatal("Can't connect to db")
	}
	queries, err := database.New(conn)
	if err != nil {
		log.Fatal()
	}

	apiCfg := apiConfig{
		DB: queries,
	}
	router := chi.NewRouter()

	router.Use(cors.Handler(cors.Options{
		AllowedOrigins: []string{"https://*", "http://*"},
		AllowedMethods: []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders: []string{"*"},
		ExposedHeaders: []string{"Link"},
		AllowCredentials: false,
		MaxAge: 300,
	}))

	v1Router := chi.NewRouter()

	v1Router.Get("/healthz", handlerReadiness)
	v1Router.Get("/err", handlerErr)
	router.Mount("/v1", v1Router)

	srv := &http.Server{
		Handler: router,
		Addr: ":" + port,
	}

	log.Printf("Server starting n port %v", port)
	err := srv.ListenAndServe()

	if err != nil {
		log.Fatal(err)
	}
}