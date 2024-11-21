package main

import (
	"encoding/json"
	"github.com/bevane/chirpy/internal/database"
	"github.com/google/uuid"
	"net/http"
	"time"
)

type Chirp struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Body      string    `json:"body"`
	UserID    uuid.UUID `json:"user_id"`
}

func (cfg *apiConfig) chirpsHandler(w http.ResponseWriter, req *http.Request) {
	type parameters struct {
		Body   string    `json:"body"`
		UserID uuid.UUID `json:"user_id"`
	}
	decoder := json.NewDecoder(req.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		errorMsg := "Couldn't decode parameters"
		respondWithError(w, http.StatusInternalServerError, errorMsg)
		return
	}
	if len(params.Body) > 140 {
		errorMsg := "Chirp is too long"
		respondWithError(w, http.StatusBadRequest, errorMsg)
		return
	}
	cleanedText := cleanProfanity(params.Body)
	chirp, err := cfg.db.CreateChirp(req.Context(), database.CreateChirpParams{
		Body:   cleanedText,
		UserID: params.UserID,
	})
	if err != nil {
		errorMsg := "Couldn't create chirp"
		respondWithError(w, http.StatusInternalServerError, errorMsg)
		return
	}
	respondWithJSON(w, http.StatusCreated, Chirp(chirp))
	return
}

func (cfg *apiConfig) getChirpsHandler(w http.ResponseWriter, req *http.Request) {
	type response struct {
		chirps []database.Chirp
	}

	chirps, err := cfg.db.GetChirps(req.Context())
	if err != nil {
		errorMsg := "Couldn't get chirps"
		respondWithError(w, http.StatusInternalServerError, errorMsg)
		return
	}
	var JSONChirps []Chirp
	for _, chirp := range chirps {
		JSONChirps = append(JSONChirps, Chirp(chirp))
	}
	respondWithJSON(w, http.StatusOK, JSONChirps)
	return
}

func (cfg *apiConfig) getChirpSingleHandler(w http.ResponseWriter, req *http.Request) {
	chirpID, err := uuid.Parse(req.PathValue("chirpID"))
	if err != nil {
		errorMsg := "Chirp not found"
		respondWithError(w, http.StatusNotFound, errorMsg)
		return
	}
	chirp, err := cfg.db.GetChirpByID(req.Context(), chirpID)
	if err != nil {
		errorMsg := "Chirp not found"
		respondWithError(w, http.StatusNotFound, errorMsg)
		return
	}
	respondWithJSON(w, http.StatusOK, Chirp(chirp))
	return
}
