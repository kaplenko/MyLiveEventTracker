package http

import (
	oauth2Models "auth-service/pkg/oauth2"
	"auth-service/pkg/oauth2/jwtUtils"
	"net/http"
)

func (h *Handler) GoogleLoginRedirect(w http.ResponseWriter, r *http.Request) {
	state, err := jwtUtils.GenerateState()
	if err != nil {
		h.log.Error("Error generating state: ", err)
	}
	authURL := h.googleService.AuthURL(state)
	http.Redirect(w, r, authURL, http.StatusTemporaryRedirect)
}

func (h *Handler) GoogleCallback(w http.ResponseWriter, r *http.Request) {
	code := r.URL.Query().Get("code")
	if code == "" {
		http.Error(w, "Missing code", http.StatusBadRequest)
		return
	}
	ctx := r.Context()
	user, err := h.googleService.UserInfo(ctx, code)
	if err != nil {
		h.log.Error("Error getting user info: ", err)
		http.Error(w, "Failed getting user info", http.StatusInternalServerError)
		return
	}

	token, err := h.useCase.AuthenticateOAuthUser(ctx, &oauth2Models.User{
		ID:           user.ID,
		Login:        user.Login,
		Email:        user.Email,
		Provider:     "google",
		AccessToken:  user.AccessToken,
		RefreshToken: user.RefreshToken,
		ExpiresAt:    user.ExpiresAt,
	})
	if err != nil {
		h.log.Error("Error authenticating user: ", err)
		http.Error(w, "Failed authenticating user", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Authorization", "Bearer "+token)

	http.SetCookie(w, &http.Cookie{
		Name:  "token",
		Value: token,
		Path:  "/",
	})
	http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
}
