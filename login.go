package main

import (
	"encoding/json"
	auth "github.com/bevane/chirpy/internal"
	"github.com/bevane/chirpy/internal/database"
	"net/http"
	"time"
)

func (cfg *apiConfig) loginHandler(w http.ResponseWriter, req *http.Request) {
	type parameters struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	type response struct {
		User
		Token        string `json:"token"`
		RefreshToken string `json:"refresh_token"`
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
	expiresInDuration := time.Hour
	token, err := auth.MakeJWT(user.ID, cfg.jwtSecret, expiresInDuration)
	if err != nil {
		errorMsg := "Error generating JWT"
		respondWithError(w, http.StatusInternalServerError, errorMsg)
		return
	}
	refreshToken, err := auth.MakeRefreshToken()
	refreshTokenExpiresAt := time.Now().Add(time.Hour * 24 * 60)
	if err != nil {
		errorMsg := "Error generating refresh token"
		respondWithError(w, http.StatusInternalServerError, errorMsg)
		return
	}
	err = cfg.db.CreateRefreshToken(req.Context(), database.CreateRefreshTokenParams{
		Token:     refreshToken,
		UserID:    user.ID,
		ExpiresAt: refreshTokenExpiresAt,
	})
	if err != nil {
		errorMsg := "Unable to add refresh token to db"
		respondWithError(w, http.StatusInternalServerError, errorMsg)
		return
	}
	respondWithJSON(w, http.StatusOK, response{
		User: User{
			ID:          user.ID,
			CreatedAt:   user.CreatedAt,
			UpdatedAt:   user.UpdatedAt,
			Email:       user.Email,
			IsChirpyRed: user.IsChirpyRed,
		},
		Token:        token,
		RefreshToken: refreshToken,
	})
	return
}

func (cfg *apiConfig) refreshHandler(w http.ResponseWriter, req *http.Request) {
	type response struct {
		Token string `json:"token"`
	}
	refreshToken, err := auth.GetBearerToken(req.Header)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	refreshTokenRow, err := cfg.db.GetRefreshToken(req.Context(), refreshToken)
	// revoked_at is only a valid time if it is not null which means the token
	// is invalid
	if err != nil || refreshTokenRow.ExpiresAt.Before(time.Now()) || refreshTokenRow.RevokedAt.Valid {
		errorMsg := "Invalid refresh token"
		respondWithError(w, http.StatusUnauthorized, errorMsg)
		return
	}
	token, err := auth.MakeJWT(refreshTokenRow.UserID, cfg.jwtSecret, time.Hour)
	respondWithJSON(w, http.StatusOK, response{
		Token: token,
	})
	return
}

func (cfg *apiConfig) revokeHandler(w http.ResponseWriter, req *http.Request) {
	type response struct {
		Token string `json:"token"`
	}
	refreshToken, err := auth.GetBearerToken(req.Header)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	err = cfg.db.RevokeRefreshToken(req.Context(), refreshToken)
	if err != nil {
		errorMsg := ("Unable to revoke refresh token")
		respondWithError(w, http.StatusInternalServerError, errorMsg)
		return
	}
	w.WriteHeader(http.StatusNoContent)
	return

}
