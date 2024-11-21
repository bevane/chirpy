package main

import (
	"database/sql"
	"fmt"
	"net/http"
	"os"
	"sync/atomic"

	"github.com/bevane/chirpy/internal/database"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

type apiConfig struct {
	fileserverHits atomic.Int32
	db             *database.Queries
	platform       string
}

func main() {
	godotenv.Load()
	dbURL := os.Getenv("DB_URL")
	platform := os.Getenv("PLATFORM")
	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		fmt.Println(err)
		return
	}
	dbQueries := database.New(db)
	apiCfg := apiConfig{
		fileserverHits: atomic.Int32{},
		db:             dbQueries,
		platform:       platform,
	}
	serveMux := http.NewServeMux()
	fileHandler := http.StripPrefix("/app/", http.FileServer(http.Dir(".")))
	serveMux.Handle("/app/", apiCfg.middlewareMetricsInc(fileHandler))
	serveMux.HandleFunc("GET /api/healthz", readyHandler)
	serveMux.HandleFunc("GET /admin/metrics", apiCfg.metricsHandler)
	serveMux.HandleFunc("POST /admin/reset", apiCfg.resetHandler)
	serveMux.HandleFunc("POST /api/validate_chirp", apiCfg.validateChirpHandler)
	serveMux.HandleFunc("POST /api/users", apiCfg.usersHandler)
	serveMux.HandleFunc("POST /api/chirps", apiCfg.chirpsHandler)
	serveMux.HandleFunc("GET /api/chirps", apiCfg.getChirpsHandler)
	server := &http.Server{
		Addr:    ":8080",
		Handler: serveMux,
	}
	server.ListenAndServe()
}
