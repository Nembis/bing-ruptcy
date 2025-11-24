package server

import "net/http"

func (s *server) handleGenerateLoginCode(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Email string
	}
}

func (s *server) handleLoginUserUsingMagicLink(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		MagicLinkCocde string
	}
}

func (s *server) handleLogOut(w http.ResponseWriter, r *http.Request) {
}
