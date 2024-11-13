package main

import (
	"encoding/json"
	"net/http"
	"strings"
)

func (cfg *apiConfig) validateChirpHandler(w http.ResponseWriter, req *http.Request) {
	type parameters struct {
		Body string `json:"body"`
	}

	decoder := json.NewDecoder(req.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		errorMsg := "Something went wrong"
		respondWithError(w, 500, errorMsg)
		return
	}
	if len(params.Body) > 140 {
		errorMsg := "Chirp is too long"
		respondWithError(w, 400, errorMsg)
		return
	}
	cleanedText := cleanProfanity(params.Body)

	type returnSuccess struct {
		CleanedBody string `json:"cleaned_body"`
	}
	respBodySuccess := returnSuccess{
		CleanedBody: cleanedText,
	}
	respondWithJSON(w, 200, respBodySuccess)
	return
}

func cleanProfanity(text string) string {
	words := strings.Split(text, " ")
	for i := range words {
		lowerWord := strings.ToLower(words[i])
		if lowerWord == "kerfuffle" || lowerWord == "sharbert" || lowerWord == "fornax" {
			words[i] = "****"
		}
	}
	return strings.Join(words, " ")
}
