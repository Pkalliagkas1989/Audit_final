package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"strings"

	"forum/config"
	"forum/models"
	"forum/repository"
	"forum/utils"
)

// AuthHandler handles authentication-related requests
type AuthHandler struct {
	UserRepo    *repository.UserRepository
	SessionRepo *repository.SessionRepository
	oauthStates map[string]string
}

// NewAuthHandler creates a new AuthHandler
func NewAuthHandler(userRepo *repository.UserRepository, sessionRepo *repository.SessionRepository) *AuthHandler {
	return &AuthHandler{
		UserRepo:    userRepo,
		SessionRepo: sessionRepo,
		oauthStates: make(map[string]string),
	}
}

// Register handles user registration
func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	// Only allow POST requests
	if r.Method != http.MethodPost {
		utils.ErrorResponse(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Parse request body
	var reg models.UserRegistration
	err := json.NewDecoder(r.Body).Decode(&reg)
	if err != nil {
		utils.ErrorResponse(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	reg.Username = strings.TrimSpace(reg.Username)
	reg.Email = strings.TrimSpace(strings.ToLower(reg.Email))
	reg.Password = strings.TrimSpace(reg.Password)

	// Validate request
	if reg.Username == "" || reg.Email == "" || reg.Password == "" {
		utils.ErrorResponse(w, "Username, email, and password are required", http.StatusBadRequest)
		return
	}

	// Username: 3–50 chars, letters/numbers/underscores only
	if !utils.UsernameRegex.MatchString(reg.Username) {
		utils.ErrorResponse(w, "Username must be 3-50 characters, letters/numbers/underscores only", http.StatusBadRequest)
		return
	}

	// Email: trim, lowercase, parse, and enforce ending in .com
	cleanEmail, err := utils.ValidateEmail(reg.Email)
	if err != nil {
		// You might want to send err.Error() directly, since ValidateEmail already produces a user-friendly message.
		utils.ErrorResponse(w, err.Error(), http.StatusBadRequest)
		return
	}
	reg.Email = cleanEmail

	// Password: at least 8 chars, at least one letter and one digit
	if !utils.IsStrongPassword(reg.Password) {
		utils.ErrorResponse(w, "Password must be at least 8 characters, with at least one letter and one digit", http.StatusBadRequest)
		return
	}

	// Optional: Strength (at least 1 digit, 1 letter)

	if !utils.IsStrongPassword(reg.Password) {
		utils.ErrorResponse(w, "Password must contain letters and numbers", http.StatusBadRequest)
		return
	}

	// Create user
	user, err := h.UserRepo.Create(reg)
	if err != nil {
		switch err {
		case repository.ErrEmailTaken:
			utils.ErrorResponse(w, "Email is already taken", http.StatusConflict)
		case repository.ErrUsernameTaken:
			utils.ErrorResponse(w, "Username is already taken", http.StatusConflict)
		default:
			utils.ErrorResponse(w, "Internal server error", http.StatusInternalServerError)
		}
		return
	}

	response := map[string]interface{}{
		"id":       user.ID,
		"username": user.Username,
		"email":    user.Email,
	}
	utils.JSONResponse(w, response, http.StatusCreated)
}

// Login handles user login
func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	// Only allow POST requests
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Parse request body
	var login models.UserLogin
	err := json.NewDecoder(r.Body).Decode(&login)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Validate request
	if login.Email == "" || login.Password == "" {
		http.Error(w, "Email and password are required", http.StatusBadRequest)
		return
	}

	// Authenticate user
	user, err := h.UserRepo.Authenticate(login)

	if err != nil {
		if err == repository.ErrInvalidCredentials {
			http.Error(w, "Invalid email or password", http.StatusUnauthorized)
		} else {
			http.Error(w, "Internal server error", http.StatusInternalServerError)
		}
		return
	}
	csrfToken := utils.GenerateCSRFToken()

	// Create a new session
	log.Println("Creating session for user:", user.ID)
	session, err := h.SessionRepo.Create(user.ID, r.RemoteAddr, csrfToken)
	log.Println("Session created:", session)
	if err != nil {
		http.Error(w, "Failed to create session", http.StatusInternalServerError)
		return
	}

	// Set cookie
	http.SetCookie(w, &http.Cookie{
		Name:     "session_id",
		Value:    session.SessionID,
		Path:     "/",
		Expires:  session.ExpiresAt,
		HttpOnly: true,
		Secure:   false,
		SameSite: http.SameSiteLaxMode, // Change from None to Lax for localhost
	})

	// Return response
	w.Header().Set("Content-Type", "application/json")
	response := models.LoginResponse{
		User:      *user,
		SessionID: session.SessionID,
		CSRFToken: csrfToken,
	}
	json.NewEncoder(w).Encode(response)
}

// Logout handles user logout
func (h *AuthHandler) Logout(w http.ResponseWriter, r *http.Request) {
	// Only allow POST requests
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Get the session cookie
	cookie, err := r.Cookie("session_id")
	if err != nil {
		// If no cookie, nothing to do
		w.WriteHeader(http.StatusOK)
		return
	}

	// Delete the session
	err = h.SessionRepo.Delete(cookie.Value)
	if err != nil {
		http.Error(w, "Failed to logout", http.StatusInternalServerError)
		return
	}

	// Clear the cookie
	http.SetCookie(w, &http.Cookie{
		Name:     "session_id",
		Value:    "",
		Path:     "/",
		MaxAge:   -1,
		HttpOnly: true,
	})

	w.WriteHeader(http.StatusOK)
}

// Inside AuthHandler
func (h *AuthHandler) VerifySession(w http.ResponseWriter, r *http.Request) {
	sessionCookie, err := r.Cookie("session_id")
	if err != nil {
		http.Error(w, "Not authenticated", http.StatusUnauthorized)
		return
	}

	session, err := h.SessionRepo.GetBySessionID(sessionCookie.Value)
	if err != nil {
		http.Error(w, "Session invalid or expired", http.StatusUnauthorized)
		return
	}

	// Optionally fetch user and return profile
	user, err := h.UserRepo.GetByID(session.UserID)
	if err != nil {
		http.Error(w, "User not found", http.StatusInternalServerError)
		return
	}

	utils.JSONResponse(w, user, http.StatusOK)
}

// GoogleLogin starts the Google OAuth flow
func (h *AuthHandler) GoogleLogin(w http.ResponseWriter, r *http.Request) {
	state := utils.GenerateCSRFToken()
	h.oauthStates[state] = "google"
	authURL := fmt.Sprintf("https://accounts.google.com/o/oauth2/v2/auth?client_id=%s&redirect_uri=%s&response_type=code&scope=openid%%20email%%20profile&state=%s",
		url.QueryEscape(config.GoogleClientID), url.QueryEscape(config.GoogleRedirectURL), state)
	http.Redirect(w, r, authURL, http.StatusFound)
}

// GoogleCallback completes the Google OAuth flow
func (h *AuthHandler) GoogleCallback(w http.ResponseWriter, r *http.Request) {
	state := r.URL.Query().Get("state")
	if h.oauthStates[state] != "google" {
		http.Error(w, "invalid state", http.StatusBadRequest)
		return
	}
	delete(h.oauthStates, state)
	code := r.URL.Query().Get("code")
	if code == "" {
		http.Error(w, "code required", http.StatusBadRequest)
		return
	}
	vals := url.Values{}
	vals.Set("client_id", config.GoogleClientID)
	vals.Set("client_secret", config.GoogleClientSecret)
	vals.Set("code", code)
	vals.Set("grant_type", "authorization_code")
	vals.Set("redirect_uri", config.GoogleRedirectURL)

	resp, err := http.PostForm("https://oauth2.googleapis.com/token", vals)
	if err != nil {
		http.Error(w, "token exchange failed", http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)
	var tok struct {
		AccessToken string `json:"access_token"`
	}
	json.Unmarshal(body, &tok)

	req, _ := http.NewRequest("GET", "https://www.googleapis.com/oauth2/v3/userinfo", nil)
	req.Header.Set("Authorization", "Bearer "+tok.AccessToken)
	uResp, err := http.DefaultClient.Do(req)
	if err != nil {
		http.Error(w, "failed userinfo", http.StatusInternalServerError)
		return
	}
	defer uResp.Body.Close()
	data, _ := io.ReadAll(uResp.Body)
	var info struct {
		Sub   string `json:"sub"`
		Email string `json:"email"`
		Name  string `json:"name"`
	}
	json.Unmarshal(data, &info)

	h.handleOAuthLogin(w, r, "google", info.Sub, info.Email, info.Name)
}

// GitHubLogin starts the GitHub OAuth flow
func (h *AuthHandler) GitHubLogin(w http.ResponseWriter, r *http.Request) {
	state := utils.GenerateCSRFToken()
	h.oauthStates[state] = "github"
	authURL := fmt.Sprintf("https://github.com/login/oauth/authorize?client_id=%s&redirect_uri=%s&scope=user:email&state=%s",
		url.QueryEscape(config.GitHubClientID), url.QueryEscape(config.GitHubRedirectURL), state)
	http.Redirect(w, r, authURL, http.StatusFound)
}

// GitHubCallback completes the GitHub OAuth flow
func (h *AuthHandler) GitHubCallback(w http.ResponseWriter, r *http.Request) {
	state := r.URL.Query().Get("state")
	if h.oauthStates[state] != "github" {
		http.Error(w, "invalid state", http.StatusBadRequest)
		return
	}
	delete(h.oauthStates, state)
	code := r.URL.Query().Get("code")
	if code == "" {
		http.Error(w, "code required", http.StatusBadRequest)
		return
	}

	vals := url.Values{}
	vals.Set("client_id", config.GitHubClientID)
	vals.Set("client_secret", config.GitHubClientSecret)
	vals.Set("code", code)
	vals.Set("redirect_uri", config.GitHubRedirectURL)

	req, _ := http.NewRequest("POST", "https://github.com/login/oauth/access_token", bytes.NewBufferString(vals.Encode()))
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		http.Error(w, "token exchange failed", http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)
	var tok struct {
		AccessToken string `json:"access_token"`
	}
	json.Unmarshal(body, &tok)

	uReq, _ := http.NewRequest("GET", "https://api.github.com/user", nil)
	uReq.Header.Set("Authorization", "Bearer "+tok.AccessToken)
	uReq.Header.Set("Accept", "application/vnd.github.v3+json")
	uResp, err := http.DefaultClient.Do(uReq)
	if err != nil {
		http.Error(w, "failed userinfo", http.StatusInternalServerError)
		return
	}
	defer uResp.Body.Close()
	data, _ := io.ReadAll(uResp.Body)
	var info struct {
		ID    int    `json:"id"`
		Email string `json:"email"`
		Login string `json:"login"`
	}
	json.Unmarshal(data, &info)

	if info.Email == "" {
		// fetch emails endpoint
		emReq, _ := http.NewRequest("GET", "https://api.github.com/user/emails", nil)
		emReq.Header.Set("Authorization", "Bearer "+tok.AccessToken)
		emReq.Header.Set("Accept", "application/vnd.github.v3+json")
		emResp, err := http.DefaultClient.Do(emReq)
		if err == nil {
			defer emResp.Body.Close()
			emData, _ := io.ReadAll(emResp.Body)
			var emails []struct {
				Email   string `json:"email"`
				Primary bool   `json:"primary"`
			}
			json.Unmarshal(emData, &emails)
			for _, e := range emails {
				if e.Primary {
					info.Email = e.Email
					break
				}
			}
		}
	}

	h.handleOAuthLogin(w, r, "github", fmt.Sprint(info.ID), info.Email, info.Login)
}

// handleOAuthLogin links or creates users based on provider info
func (h *AuthHandler) handleOAuthLogin(w http.ResponseWriter, r *http.Request, provider, providerID, email, name string) {
	if email == "" {
		http.Error(w, "email not provided", http.StatusBadRequest)
		return
	}

	user, err := h.UserRepo.GetByProvider(provider, providerID)
	if err == repository.ErrUserNotFound {
		existing, err2 := h.UserRepo.GetByEmail(strings.ToLower(email))
		if err2 == nil {
			h.UserRepo.LinkProvider(existing.ID, provider, providerID)
			user = existing
		} else if err2 == repository.ErrUserNotFound {
			username := name
			if username == "" {
				username = strings.Split(email, "@")[0]
			}
			user, err = h.UserRepo.CreateWithProvider(username, strings.ToLower(email), provider, providerID)
			if err != nil {
				http.Error(w, "failed to create user", http.StatusInternalServerError)
				return
			}
		} else {
			http.Error(w, "internal error", http.StatusInternalServerError)
			return
		}
	} else if err != nil {
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}

	csrf := utils.GenerateCSRFToken()
	session, err := h.SessionRepo.Create(user.ID, r.RemoteAddr, csrf)
	if err != nil {
		http.Error(w, "failed to create session", http.StatusInternalServerError)
		return
	}
	http.SetCookie(w, &http.Cookie{
		Name:     "session_id",
		Value:    session.SessionID,
		Path:     "/",
		Expires:  session.ExpiresAt,
		HttpOnly: true,
		Secure:   false,
		SameSite: http.SameSiteLaxMode,
	})
	resp := models.LoginResponse{User: *user, SessionID: session.SessionID, CSRFToken: csrf}
	json.NewEncoder(w).Encode(resp)
}
