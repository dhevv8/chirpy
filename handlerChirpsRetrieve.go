package main 

import (
	"net/http"
)

func (cfg *apiConfig) handlerChirpsRetrieve(w http.ResponseWriter, r *http.Request) {
	dbChirps, err := cfg.db.GetChirp(r.Context())
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't retrieve chirps", err)
		return
	}
	chirps := []Chirp{}
	for _, dbChirp := range dbChirps {
		chirps = append(chirps, Chirp{
			Id:        dbChirp.ID,
			Created_at: dbChirp.CreatedAt,
			Updated_at: dbChirp.UpdatedAt,
			Body:      dbChirp.Body,
			User_id:    dbChirp.UserID,
		})
	}
	respondWithJSON(w, http.StatusOK, chirps)
}