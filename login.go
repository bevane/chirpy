package main

import (
	"encoding/json"
	"net/http"
	"time"

	auth "github.com/bevane/chirpy/internal"
)

func (cfg *apiConfig) loginHandler(w http.ResponseWriter, req *http.Request) {
	type parameters struct {
		Email            string `json:"email"`
		Password         string `json:"password"`
		ExpiresInSeconds int    `json:"expires_in_seconds"`
	}
	type response struct {
		User
		Token string `json:"token"`
	}
	decoder := json.NewDecoder(req.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		errorMsg := "Couldn't decode parameters"
		respondWithError(w, http.StatusInternalServerError, errorMsg)
		return
	}
	if params.Password == "" || params.Email == "" {
		errorMsg := "Password and Email cannot be empty"
		respondWithError(w, http.StatusBadRequest, errorMsg)
		return

	}
	expiresInDuration := time.Duration(params.ExpiresInSeconds) * time.Second
	if expiresInDuration == 0 || expiresInDuration > time.Hour {
		expiresInDuration = time.Hour
	}
	user, err := cfg.db.GetUserByEmail(req.Context(), params.Email)
	if err != nil {
		errorMsg := "User not found"
		respondWithError(w, http.StatusInternalServerError, errorMsg)
		return
	}
	err = auth.CheckPasswordHash(params.Password, user.HashedPassword)
	if err != nil {
		errorMsg := "Incorrect password"
		respondWithError(w, http.StatusUnauthorized, errorMsg)
		return
	}
	token, err := auth.MakeJWT(user.ID, cfg.jwtSecret, expiresInDuration)
	if err != nil {
		errorMsg := "Error generating JWT"
		respondWithError(w, http.StatusInternalServerError, errorMsg)
		return
	}
	respondWithJSON(w, http.StatusOK, response{
		User: User{
			ID:        user.ID,
			CreatedAt: user.CreatedAt,
			UpdatedAt: user.UpdatedAt,
			Email:     user.Email,
		},
		Token: token,
	})
	return
}
