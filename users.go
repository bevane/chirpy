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
	ID          uuid.UUID `json:"id"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	Email       string    `json:"email"`
	IsChirpyRed bool      `json:"is_chirpy_red"`
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

func (cfg *apiConfig) updateUserHandler(w http.ResponseWriter, req *http.Request) {
	type parameters struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	type response struct {
		User
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
	user, err := cfg.db.UpdateUserEmailAndPassword(req.Context(), database.UpdateUserEmailAndPasswordParams{
		Email:          params.Email,
		HashedPassword: hashedPassword,
		ID:             userID,
	})
	if err != nil {
		errorMsg := "Couldn't update user"
		respondWithError(w, http.StatusInternalServerError, errorMsg)
		return
	}

	respondWithJSON(w, http.StatusOK, response{
		User: User(user),
	})
	return
}
