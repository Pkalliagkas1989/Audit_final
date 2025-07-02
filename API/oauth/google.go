package oauth

import (
	"encoding/json"
	"net/http"
	"net/url"
)

// GoogleProvider implements OAuth for Google.
type GoogleProvider struct {
	ClientID     string
	ClientSecret string
	RedirectURL  string
}

func (g *GoogleProvider) Name() string { return "google" }

func (g *GoogleProvider) LoginURL(state string) string {
	v := url.Values{}
	v.Set("client_id", g.ClientID)
	v.Set("redirect_uri", g.RedirectURL)
	v.Set("response_type", "code")
	v.Set("scope", "email")
	v.Set("state", state)
	return "https://accounts.google.com/o/oauth2/v2/auth?" + v.Encode()
}

func (g *GoogleProvider) Exchange(code string) (*UserInfo, error) {
	data := url.Values{}
	data.Set("client_id", g.ClientID)
	data.Set("client_secret", g.ClientSecret)
	data.Set("code", code)
	data.Set("grant_type", "authorization_code")
	data.Set("redirect_uri", g.RedirectURL)
	resp, err := http.PostForm("https://oauth2.googleapis.com/token", data)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	var token struct {
		AccessToken string `json:"access_token"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&token); err != nil {
		return nil, err
	}
	req, _ := http.NewRequest("GET", "https://www.googleapis.com/oauth2/v2/userinfo", nil)
	req.Header.Set("Authorization", "Bearer "+token.AccessToken)
	uiResp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer uiResp.Body.Close()
	var info struct {
		ID    string `json:"id"`
		Email string `json:"email"`
	}
	if err := json.NewDecoder(uiResp.Body).Decode(&info); err != nil {
		return nil, err
	}
	return &UserInfo{ID: info.ID, Email: info.Email}, nil
}
