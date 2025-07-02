package handlers

import (
	"encoding/json"
	"net/http"
	"os"

	"forum/models"
	"forum/oauth"
	"forum/repository"
	"forum/utils"
)

// OAuthHandler manages OAuth logins
type OAuthHandler struct {
	userRepo     *repository.UserRepository
	providerRepo *repository.UserProviderRepository
	sessionRepo  *repository.SessionRepository
	providers    map[string]oauth.Provider
}

// NewOAuthHandler creates a new handler
func NewOAuthHandler(userRepo *repository.UserRepository, providerRepo *repository.UserProviderRepository, sessionRepo *repository.SessionRepository) *OAuthHandler {
	providers := map[string]oauth.Provider{
		"google": &oauth.GoogleProvider{
			ClientID:     os.Getenv("GOOGLE_CLIENT_ID"),
			ClientSecret: os.Getenv("GOOGLE_CLIENT_SECRET"),
			RedirectURL:  "http://localhost:8080/forum/api/oauth/google/callback",
		},
		"github": &oauth.GithubProvider{
			ClientID:     os.Getenv("GITHUB_CLIENT_ID"),
			ClientSecret: os.Getenv("GITHUB_CLIENT_SECRET"),
			RedirectURL:  "http://localhost:8080/forum/api/oauth/github/callback",
		},
	}
	return &OAuthHandler{userRepo: userRepo, providerRepo: providerRepo, sessionRepo: sessionRepo, providers: providers}
}

func (h *OAuthHandler) login(providerName string, w http.ResponseWriter, r *http.Request) {
	p, ok := h.providers[providerName]
	if !ok {
		http.Error(w, "unknown provider", http.StatusBadRequest)
		return
	}
	state := utils.GenerateState()
	http.Redirect(w, r, p.LoginURL(state), http.StatusTemporaryRedirect)
}

// GoogleLogin redirects to Google OAuth
func (h *OAuthHandler) GoogleLogin(w http.ResponseWriter, r *http.Request) {
	h.login("google", w, r)
}

// GithubLogin redirects to GitHub OAuth
func (h *OAuthHandler) GithubLogin(w http.ResponseWriter, r *http.Request) {
	h.login("github", w, r)
}

func (h *OAuthHandler) handleCallback(w http.ResponseWriter, r *http.Request, providerName string) {
	code := r.URL.Query().Get("code")
	if code == "" {
		http.Error(w, "missing code", http.StatusBadRequest)
		return
	}
	provider, ok := h.providers[providerName]
	if !ok {
		http.Error(w, "unknown provider", http.StatusBadRequest)
		return
	}
	info, err := provider.Exchange(code)
	if err != nil {
		http.Error(w, "failed to exchange code", http.StatusBadRequest)
		return
	}
	user, err := h.userRepo.GetOrCreateOAuthUser(info.Email, "", providerName, info.ID)
	if err != nil {
		utils.ErrorResponse(w, "could not process user", http.StatusInternalServerError)
		return
	}
	csrfToken := utils.GenerateCSRFToken()
	session, err := h.sessionRepo.Create(user.ID, r.RemoteAddr, csrfToken)
	if err != nil {
		http.Error(w, "failed to create session", http.StatusInternalServerError)
		return
	}
	http.SetCookie(w, &http.Cookie{Name: "session_id", Value: session.SessionID, Path: "/", Expires: session.ExpiresAt, HttpOnly: true})
	utils.JSONResponse(w, models.LoginResponse{User: *user, SessionID: session.SessionID, CSRFToken: csrfToken}, http.StatusOK)
}

// GoogleCallback handles Google OAuth callback
func (h *OAuthHandler) GoogleCallback(w http.ResponseWriter, r *http.Request) {
	h.handleCallback(w, r, "google")
}

// GithubCallback handles GitHub OAuth callback
func (h *OAuthHandler) GithubCallback(w http.ResponseWriter, r *http.Request) {
	h.handleCallback(w, r, "github")
}
