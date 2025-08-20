package main

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"
	"time"

	"github.com/dhevv8/chirpy/internal/database"
	"github.com/google/uuid"
)
type Chirp struct{
		Id uuid.UUID `json:"id"`
		Created_at time.Time `json:"created_at"`
		Updated_at time.Time  `json:"updated_at"`
		Body string `json:"body"`
		User_id uuid.UUID `json:"user_id"`
	}

func (cfg *apiConfig) handleChirps(w http.ResponseWriter,r *http.Request){
	type parameters struct{
		Body string `json:"body"`
		User_id uuid.UUID `json:"user_id"`
	}

	decoder:= json.NewDecoder(r.Body)
	params:=parameters{}
	err:=decoder.Decode(&params)
	if err!=nil{
		respondWithError(w,http.StatusInternalServerError,"Could not decode request body",err)
		return
	}
	cleaned,err:=validateChirp(params.Body)
	if err!=nil{
		respondWithError(w,http.StatusBadRequest,"Invalid chirp body",err)
		return
	}
	chirp,err:=cfg.db.CreateChirp(r.Context(),database.CreateChirpParams{
		Body: cleaned,
		UserID: params.User_id,
	})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't create chirp", err)
		return
	}
	respondWithJSON(w, http.StatusCreated, Chirp{
		Id:        chirp.ID,
		Created_at: chirp.CreatedAt,
		Updated_at: chirp.UpdatedAt,
		Body:      chirp.Body,
		User_id:    chirp.UserID,
	})
}
func validateChirp(body string) (string, error) {
	const maxChirpLength = 140
	if len(body) > maxChirpLength {
		return "", errors.New("Chirp is too long")
	}

	badWords := map[string]struct{}{
		"kerfuffle": {},
		"sharbert":  {},
		"fornax":    {},
	}
	cleaned := getCleanedBody(body, badWords)
	return cleaned, nil
}

func getCleanedBody(body string, badWords map[string]struct{}) string {
	words:=strings.Split(body," ")
	for i,word:=range words{
		loweredWord:=strings.ToLower(word)
		if _,ok:=badWords[loweredWord];ok{
			words[i]="****"
		}
	}
	cleaned:=strings.Join(words," ")
	return cleaned
}