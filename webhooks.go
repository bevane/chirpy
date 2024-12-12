package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	auth "github.com/bevane/chirpy/internal"
	"github.com/google/uuid"
)

func (cfg *apiConfig) polkaWebhookHandler(w http.ResponseWriter, req *http.Request) {
	type parameters struct {
		Event string `json:"event"`
		Data  struct {
			UserID uuid.UUID `json:"user_id"`
		} `json:"data"`
	}
	type response struct {
		User
	}
	key, err := auth.GetAPIKey(req.Header)
	if err != nil || key != cfg.polkaKey {
		errorMsg := fmt.Sprintf("Invalid api key %s", err.Error())
		respondWithError(w, http.StatusUnauthorized, errorMsg)
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
	if params.Event != "user.upgraded" {
		w.WriteHeader(http.StatusNoContent)
		return
	}
	err = cfg.db.UpgradeUserByID(req.Context(), params.Data.UserID)
	if err != nil {
		errorMsg := "User not found"
		respondWithError(w, http.StatusNotFound, errorMsg)
		return
	}
	w.WriteHeader(http.StatusNoContent)
	return
}
