package main

import (
	"encoding/json"
	auth "github.com/bevane/chirpy/internal"
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
		Body string `json:"body"`
	}
	token, err := auth.GetBearerToken(req.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, err.Error())
		return
	}
	userID, err := auth.ValidateJWT(token, cfg.jwtSecret)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, err.Error())
		return
	}
	decoder := json.NewDecoder(req.Body)
	params := parameters{}
	err = decoder.Decode(&params)
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
		UserID: userID,
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
	var chirps []database.Chirp
	var reverseSort bool
	var err error
	if sortStr := req.URL.Query().Get("sort"); sortStr == "desc" {
		reverseSort = true
	} else if sortStr == "" || sortStr == "asc" {
		reverseSort = false
	} else {
		errorMsg := "sort paramter not valid"
		respondWithError(w, http.StatusBadRequest, errorMsg)
		return
	}
	if userIDstr := req.URL.Query().Get("author_id"); userIDstr != "" {
		userID, err := uuid.Parse(userIDstr)
		if err != nil {
			errorMsg := "author_id not a valid uuid"
			respondWithError(w, http.StatusBadRequest, errorMsg)
			return
		}
		chirps, err = cfg.db.GetChirpsByUserID(req.Context(), userID)
		if err != nil {
			errorMsg := "unable to chirps by author"
			respondWithError(w, http.StatusInternalServerError, errorMsg)
			return
		}
	} else {
		chirps, err = cfg.db.GetChirps(req.Context())
		if err != nil {
			errorMsg := "Couldn't get chirps"
			respondWithError(w, http.StatusInternalServerError, errorMsg)
			return
		}
	}
	var JSONChirps []Chirp
	if reverseSort {
		for i := len(chirps) - 1; i >= 0; i-- {
			JSONChirps = append(JSONChirps, Chirp(chirps[i]))
		}
	} else {
		for _, chirp := range chirps {
			JSONChirps = append(JSONChirps, Chirp(chirp))
		}
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

func (cfg *apiConfig) deleteChirpSingleHandler(w http.ResponseWriter, req *http.Request) {
	chirpID, err := uuid.Parse(req.PathValue("chirpID"))
	if err != nil {
		errorMsg := "Chirp id not valid"
		respondWithError(w, http.StatusNotFound, errorMsg)
		return
	}
	token, err := auth.GetBearerToken(req.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, err.Error())
		return
	}
	userID, err := auth.ValidateJWT(token, cfg.jwtSecret)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, err.Error())
		return
	}
	chirp, err := cfg.db.GetChirpByID(req.Context(), chirpID)
	if err != nil {
		errorMsg := "Chirp not found"
		respondWithError(w, http.StatusNotFound, errorMsg)
		return
	}
	if chirp.UserID != userID {
		errorMsg := "You are not authorized to delete that chirp"
		respondWithError(w, http.StatusForbidden, errorMsg)
		return
	}
	err = cfg.db.DeleteChirpByID(req.Context(), chirpID)
	if err != nil {
		errorMsg := "Unable to delete chirp"
		respondWithError(w, http.StatusInternalServerError, errorMsg)
		return
	}
	w.WriteHeader(http.StatusNoContent)
	return
}
