package http

import (
	"auth-service/pkg/oauth2/github"
	"fmt"
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
	ctx := r.Context()
	user, err := h.githubService.UserInfo(ctx, code)
	if err != nil {
		h.log.Error("Error getting user info: ", err)
		http.Error(w, "Failed to get user info", http.StatusInternalServerError)
		return
	}
	token, err := h.usecase.AuthenticateOAuthUser(ctx, &github.User{
		ID:           user.ID,
		Login:        user.Name,
		Name:         user.Name,
		Email:        user.Email,
		Provider:     "github",
		AccessToken:  user.AccessToken,
		RefreshToken: user.RefreshToken,
		ExpiresAt:    user.ExpiresAt,
	})
	if err != nil {
		h.log.Error("Error authenticating user", "error", err)
		http.Error(w, "Authentication failed", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Authorization", "Bearer "+token)
	w.Write([]byte(fmt.Sprintf("Welcome, %s!", user.Name)))
}
