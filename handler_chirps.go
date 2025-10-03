package main

import (
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"github.com/geolunalg/gochirpy/internal/auth"
	"github.com/geolunalg/gochirpy/internal/database"
	"github.com/google/uuid"
)

type returnVals struct {
	ID        string    `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Body      string    `json:"body"`
	UserID    string    `json:"user_id"`
}

func (cfg *apiConfig) handlerAddChirp(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Body   string `json:"body"`
		// UserID string `json:"user_id"`
	}

	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "invalid token", err)
		return
	}

	userId, err := auth.ValidateJWT(token, cfg.jwtSecret)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "invalid token", err)
		return
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err = decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't decode parameters", err)
		return
	}

	const maxChirpLength = 140
	if len(params.Body) > maxChirpLength {
		respondWithError(w, http.StatusBadRequest, "Chirp is too long", nil)
		return
	}

	badWords := map[string]struct{}{
		"kerfuffle": {},
		"sharbert":  {},
		"fornax":    {},
	}
	cleaned := getCleanedBody(params.Body, badWords)

	// userUUID, err := uuid.Parse(params.UserID)
	// if err != nil {
	// 	respondWithError(w, http.StatusInternalServerError, "Couldn't parse user id", err)
	// 	return
	// }

	chirp, err := cfg.db.CreateChirp(r.Context(), database.CreateChirpParams{
		ID:        uuid.New(),
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		Body:      cleaned,
		UserID:    userId,
	})
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Failed to create new chirp", err)
		return
	}

	resp := returnVals{
		ID:        chirp.ID.String(),
		CreatedAt: chirp.CreatedAt,
		UpdatedAt: chirp.UpdatedAt,
		Body:      chirp.Body,
		UserID:    chirp.UserID.String(),
	}
	respondWithJSON(w, http.StatusCreated, resp)
}

func getCleanedBody(body string, badWords map[string]struct{}) string {
	words := strings.Split(body, " ")
	for i, word := range words {
		loweredWord := strings.ToLower(word)
		if _, ok := badWords[loweredWord]; ok {
			words[i] = "****"
		}
	}
	cleaned := strings.Join(words, " ")
	return cleaned
}

func (cfg *apiConfig) handlerGetChirps(w http.ResponseWriter, r *http.Request) {
	allChirps := []returnVals{}

	chirps, err := cfg.db.GetChirps(r.Context())
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Failed to get chirp", err)
		return
	}

	for _, chirp := range chirps {
		respChirp := returnVals{
			ID:        chirp.ID.String(),
			CreatedAt: chirp.CreatedAt,
			UpdatedAt: chirp.UpdatedAt,
			Body:      chirp.Body,
			UserID:    chirp.UserID.String(),
		}
		allChirps = append(allChirps, respChirp)
	}

	respondWithJSON(w, http.StatusOK, allChirps)
}

func (cfg *apiConfig) handlerGetChirpById(w http.ResponseWriter, r *http.Request) {
	chirpId := r.PathValue("chirpID")
	chirpUUID, err := uuid.Parse(chirpId)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't parse chirp id", err)
		return
	}

	chirp, err := cfg.db.GetChirpById(r.Context(), chirpUUID)
	if err != nil {
		respondWithError(w, http.StatusNotFound, "Failed to get chirp by id", err)
		return
	}

	resp := returnVals{
		ID:        chirp.ID.String(),
		CreatedAt: chirp.CreatedAt,
		UpdatedAt: chirp.UpdatedAt,
		Body:      chirp.Body,
		UserID:    chirp.UserID.String(),
	}
	respondWithJSON(w, http.StatusOK, resp)
}
