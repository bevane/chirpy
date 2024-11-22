package main

import (
	"encoding/json"
	"net/http"
	"time"

	auth "github.com/bevane/chirpy/internal"
	"github.com/bevane/chirpy/internal/database"
	"github.com/google/uuid"
)

type User struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Email     string    `json:"email"`
}

func (cfg *apiConfig) usersHandler(w http.ResponseWriter, req *http.Request) {
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
	if params.Password == "" {
		errorMsg := "Password cannot be empty"
		respondWithError(w, http.StatusBadRequest, errorMsg)
		return

	}
	hashedPassword, err := auth.HashPassword(params.Password)
	if err != nil {
		errorMsg := "Couldn't hash password"
		respondWithError(w, http.StatusInternalServerError, errorMsg)
		return
	}
	user, err := cfg.db.CreateUser(req.Context(), database.CreateUserParams{
		Email:          params.Email,
		HashedPassword: hashedPassword,
	})
	if err != nil {
		errorMsg := "Couldn't create user"
		respondWithError(w, http.StatusInternalServerError, errorMsg)
		return
	}

	respondWithJSON(w, http.StatusCreated, response{
		User: User(user),
	})
	return
}
