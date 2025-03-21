package http

import (
	"auth-service/pkg/oauth2/github"
	"errors"
	"github.com/jackc/pgx/v5"
	"net/http"
)

func (h *Handler) GithubLoginRedirect(w http.ResponseWriter, r *http.Request) {
	state, err := github.GenerateState()
	if err != nil {
		h.log.Error("Error generating state: ", err)
	}
	authURL := h.githubService.AuthURL(state)
	http.Redirect(w, r, authURL, http.StatusFound)
}

func (h *Handler) GithubCallback(w http.ResponseWriter, r *http.Request) {
	code := r.URL.Query().Get("code")
	if code == "" {
		http.Error(w, "Missing code", http.StatusBadRequest)
		return
	}
	user, err := h.githubService.UserInfo(r.Context(), code)
	if err != nil {
		h.log.Error("Error getting user info: ", err)
	}

	existingUser, err := h.usecase.GetUser(r.Context(), user.ID)
	if err != nil && errors.Is(err, pgx.ErrNoRows) {
		http.Error(w, "Failed to query the database", http.StatusInternalServerError)
		return
	}

	if existingUser != nil {

	}
}
