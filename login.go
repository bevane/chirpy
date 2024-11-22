package main

import (
	"encoding/json"
	auth "github.com/bevane/chirpy/internal"
	"net/http"
)

func (cfg *apiConfig) loginHandler(w http.ResponseWriter, req *http.Request) {
	type parameters struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	type response struct {
		User
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

	respondWithJSON(w, http.StatusOK, User{
		ID:        user.ID,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
		Email:     user.Email,
	})
	return
}
