package main

import (
	"encoding/json"
	"net/http"
	"time"
	"github.com/google/uuid"
	"github.com/dhevv8/chirpy/internal/database"
	"github.com/dhevv8/chirpy/internal/auth"
)

type User struct{
	ID  uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Email     string    `json:"email"`
	Password  string    `json:"-"`
}

func (cfg *apiConfig) handlerUsers(w http.ResponseWriter, r *http.Request){
	type parameters struct{
		Password string `json:"password"`
		Email string `json:"email"`
	}
	type response struct{
		User
	}

	decoder:=json.NewDecoder(r.Body)
	params:=parameters{}
	err:=decoder.Decode(&params)
	if err!=nil{
		respondWithError(w,http.StatusInternalServerError, "Invalid request body",err)
		return
	}
	hashedPassword,err:=auth.HashPassword(params.Password)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't hash password", err)
		return
	}
	user, err := cfg.db.CreateUser(r.Context(), database.CreateUserParams{
		Email:          params.Email,
		HashedPassword: hashedPassword,
	})
	if err!=nil{
		respondWithError(w,http.StatusInternalServerError, "Failed to create user",err)
		return
	}
	respondWithJSON(w,http.StatusCreated,response{User:User{
		ID:        user.ID,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
		Email:    user.Email,
	}})
}