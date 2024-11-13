package main

import (
	"encoding/json"
	"net/http"
)

func respondWithError(w http.ResponseWriter, code int, msg string) {
	type returnError struct {
		Error string `json:"error"`
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	respBodyError := returnError{
		Error: msg,
	}
	dat, _ := json.Marshal(respBodyError)
	w.Write(dat)
	return

}

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	dat, _ := json.Marshal(payload)
	w.Write(dat)
	return

}
