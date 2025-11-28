package server

import (
	"crypto/rand"
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/nembis/bing-ruptcy/backend/internal/database"
)

func (s *server) handleGenerateLoginCode(w http.ResponseWriter, r *http.Request) {
	logger := getEnrichedLogger(r.Context())
	type parameters struct {
		Email string
	}
	data := &parameters{}

	dec := json.NewDecoder(r.Body)

	if err := dec.Decode(data); err != nil {
		logger.ErrorContext(r.Context(), "Failed to decode request body as json", "Error", err)
		respondWithError(w, http.StatusBadRequest, "Missing body")
	}

	magicLink := make([]byte, 32)

	_, err := rand.Read(magicLink)
	if err != nil {
		logger.ErrorContext(r.Context(), "Failed to generate magic link token")
	}

	magicLinkPG := pgtype.Text{
		String: string(magicLink),
		Valid:  true,
	}

	params := database.InsertMagicLinkParams{
		Email:     data.Email,
		MagicLink: magicLinkPG,
	}

	_, err = s.queries.InsertMagicLink(r.Context(), params)
	if errors.Is(err, sql.ErrNoRows) {
		if _, err = s.queries.CreateUser(r.Context(), database.CreateUserParams{
			UserName:  data.Email,
			Email:     data.Email,
			MagicLink: magicLinkPG,
		},
		); err != nil {
			logger.ErrorContext(r.Context(), "Failed to create new user")
			respondWithError(w, http.StatusInternalServerError, "Failed to create new user")
		}
		logger.InfoContext(r.Context(), "Created user with email due to new user requesting magic link")
	}
	// TODO: Add emailing of the link. For now sending the magic link in the json.

	respondWithJSON(w, http.StatusOK, magicLink)
}

func (s *server) handleLoginUserUsingMagicLink(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		MagicLinkCode string
	}

}

func (s *server) handleLogOut(w http.ResponseWriter, r *http.Request) {

}
