package main

import "net/http"

func (cfg *apiConfig) resetHandler(w http.ResponseWriter, req *http.Request) {
	if cfg.platform != "dev" {
		w.WriteHeader(http.StatusForbidden)
		w.Write([]byte("Reset is only allowed in dev environment."))
		return
	}
	cfg.fileserverHits.Store(0)

	err := cfg.db.DeleteAllUsers(req.Context())
	if err != nil {
		errorMsg := "Error deleting users"
		respondWithError(w, 500, errorMsg)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Hits and users database reset successful"))
}
