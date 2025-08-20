package main

import(
	"net/http"
	"github.com/google/uuid"
)

type ChirpByID struct{
	Body string `json:"body"`
}

func (cfg *apiConfig) handlerChirpByID(w http.ResponseWriter, r *http.Request) {
	chirpIDStr := r.PathValue("chirpID")
	uuidChirpID, err := uuid.Parse(chirpIDStr)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid chirp ID format", err)
		return
	}
	chirp, err := cfg.db.GetOneChirp(r.Context(), uuidChirpID)
	if err != nil {
		respondWithError(w, http.StatusNotFound, "Couldn't retrieve chirp", err)
		return
	}

	respondWithJSON(w, http.StatusOK, ChirpByID{
		Body:chirp,
	})
}
