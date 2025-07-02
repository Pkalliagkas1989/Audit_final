package oauth

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
)

// GithubProvider implements OAuth for GitHub.
type GithubProvider struct {
	ClientID     string
	ClientSecret string
	RedirectURL  string
}

func (g *GithubProvider) Name() string { return "github" }

func (g *GithubProvider) LoginURL(state string) string {
	v := url.Values{}
	v.Set("client_id", g.ClientID)
	v.Set("redirect_uri", g.RedirectURL)
	v.Set("scope", "user:email")
	v.Set("state", state)
	return "https://github.com/login/oauth/authorize?" + v.Encode()
}

func (g *GithubProvider) Exchange(code string) (*UserInfo, error) {
	data := url.Values{}
	data.Set("client_id", g.ClientID)
	data.Set("client_secret", g.ClientSecret)
	data.Set("code", code)
	data.Set("redirect_uri", g.RedirectURL)

	req, _ := http.NewRequest("POST", "https://github.com/login/oauth/access_token", strings.NewReader(data.Encode()))
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := http.DefaultClient.Do(req)
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
	userReq, _ := http.NewRequest("GET", "https://api.github.com/user", nil)
	userReq.Header.Set("Authorization", "Bearer "+token.AccessToken)
	uiResp, err := http.DefaultClient.Do(userReq)
	if err != nil {
		return nil, err
	}
	defer uiResp.Body.Close()
	var base struct {
		ID    int64  `json:"id"`
		Email string `json:"email"`
	}
	if err := json.NewDecoder(uiResp.Body).Decode(&base); err != nil {
		return nil, err
	}
	email := base.Email
	if email == "" {
		emailsReq, _ := http.NewRequest("GET", "https://api.github.com/user/emails", nil)
		emailsReq.Header.Set("Authorization", "Bearer "+token.AccessToken)
		emailsResp, err := http.DefaultClient.Do(emailsReq)
		if err != nil {
			return nil, err
		}
		defer emailsResp.Body.Close()
		var emails []struct {
			Email    string `json:"email"`
			Primary  bool   `json:"primary"`
			Verified bool   `json:"verified"`
		}
		if err := json.NewDecoder(emailsResp.Body).Decode(&emails); err != nil {
			return nil, err
		}
		for _, e := range emails {
			if e.Primary && e.Verified {
				email = e.Email
				break
			}
		}
	}
	return &UserInfo{ID: fmt.Sprintf("%d", base.ID), Email: email}, nil
}
