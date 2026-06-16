package usecase

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
)

type OAuthConfig struct {
	GoogleID       string
	GoogleSecret   string
	GitHubID       string
	GitHubSecret   string
	FrontendURL    string
	CallbackBase   string
}

const oauthStateCookie = "oauth_state"

func googleRedirectURL(cfg OAuthConfig, state string) string {
	v := url.Values{
		"client_id":    {cfg.GoogleID},
		"redirect_uri": {cfg.CallbackBase + "/google"},
		"response_type": {"code"},
		"scope":         {"openid email profile"},
		"state":         {state},
	}
	return "https://accounts.google.com/o/oauth2/v2/auth?" + v.Encode()
}

func githubRedirectURL(cfg OAuthConfig, state string) string {
	v := url.Values{
		"client_id":    {cfg.GitHubID},
		"redirect_uri": {cfg.CallbackBase + "/github"},
		"scope":        {"read:user user:email"},
		"state":        {state},
	}
	return "https://github.com/login/oauth/authorize?" + v.Encode()
}

type googleTokenResponse struct {
	AccessToken string `json:"access_token"`
	IDToken     string `json:"id_token"`
}

type googleUserInfo struct {
	Sub   string `json:"sub"`
	Email string `json:"email"`
	Name  string `json:"name"`
}

func exchangeGoogleCode(ctx context.Context, cfg OAuthConfig, code string) (*googleUserInfo, error) {
	v := url.Values{
		"code":          {code},
		"client_id":     {cfg.GoogleID},
		"client_secret": {cfg.GoogleSecret},
		"redirect_uri":  {cfg.CallbackBase + "/google"},
		"grant_type":    {"authorization_code"},
	}
	req, err := http.NewRequestWithContext(ctx, "POST", "https://oauth2.googleapis.com/token", strings.NewReader(v.Encode()))
	if err != nil {
		return nil, fmt.Errorf("google token request: %w", err)
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("google token exchange: %w", err)
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("google token read: %w", err)
	}
	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("google token error (HTTP %d): %s", resp.StatusCode, string(body))
	}
	var tr googleTokenResponse
	if err := json.Unmarshal(body, &tr); err != nil {
		return nil, fmt.Errorf("google token decode: %w", err)
	}
	req2, err := http.NewRequestWithContext(ctx, "GET", "https://www.googleapis.com/oauth2/v2/userinfo", nil)
	if err != nil {
		return nil, fmt.Errorf("google userinfo request: %w", err)
	}
	req2.Header.Set("Authorization", "Bearer "+tr.AccessToken)
	resp2, err := http.DefaultClient.Do(req2)
	if err != nil {
		return nil, fmt.Errorf("google userinfo: %w", err)
	}
	defer resp2.Body.Close()
	body2, err := io.ReadAll(resp2.Body)
	if err != nil {
		return nil, fmt.Errorf("google userinfo read: %w", err)
	}
	if resp2.StatusCode != 200 {
		return nil, fmt.Errorf("google userinfo error (HTTP %d): %s", resp2.StatusCode, string(body2))
	}
	var ui googleUserInfo
	if err := json.Unmarshal(body2, &ui); err != nil {
		return nil, fmt.Errorf("google userinfo decode: %w", err)
	}
	return &ui, nil
}

type githubTokenResponse struct {
	AccessToken string `json:"access_token"`
}

type githubUserInfo struct {
	ID    int64  `json:"id"`
	Email string `json:"email"`
	Name  string `json:"name"`
	Login string `json:"login"`
}

type githubEmailInfo struct {
	Email    string `json:"email"`
	Primary  bool   `json:"primary"`
	Verified bool   `json:"verified"`
}

func exchangeGitHubCode(ctx context.Context, cfg OAuthConfig, code string) (*githubUserInfo, error) {
	v := url.Values{
		"code":          {code},
		"client_id":     {cfg.GitHubID},
		"client_secret": {cfg.GitHubSecret},
		"redirect_uri":  {cfg.CallbackBase + "/github"},
	}
	req, err := http.NewRequestWithContext(ctx, "POST", "https://github.com/login/oauth/access_token", strings.NewReader(v.Encode()))
	if err != nil {
		return nil, fmt.Errorf("github token request: %w", err)
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Accept", "application/json")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("github token exchange: %w", err)
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("github token read: %w", err)
	}
	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("github token error (HTTP %d): %s", resp.StatusCode, string(body))
	}
	var tr githubTokenResponse
	if err := json.Unmarshal(body, &tr); err != nil {
		return nil, fmt.Errorf("github token decode: %w", err)
	}
	req2, err := http.NewRequestWithContext(ctx, "GET", "https://api.github.com/user", nil)
	if err != nil {
		return nil, fmt.Errorf("github userinfo request: %w", err)
	}
	req2.Header.Set("Authorization", "Bearer "+tr.AccessToken)
	req2.Header.Set("Accept", "application/json")
	resp2, err := http.DefaultClient.Do(req2)
	if err != nil {
		return nil, fmt.Errorf("github userinfo: %w", err)
	}
	defer resp2.Body.Close()
	body2, err := io.ReadAll(resp2.Body)
	if err != nil {
		return nil, fmt.Errorf("github userinfo read: %w", err)
	}
	if resp2.StatusCode != 200 {
		return nil, fmt.Errorf("github userinfo error (HTTP %d): %s", resp2.StatusCode, string(body2))
	}
	var ui githubUserInfo
	if err := json.Unmarshal(body2, &ui); err != nil {
		return nil, fmt.Errorf("github userinfo decode: %w", err)
	}
	if ui.Email == "" {
		emails, err := fetchGitHubEmails(ctx, tr.AccessToken)
		if err == nil {
			for _, e := range emails {
				if e.Primary && e.Verified {
					ui.Email = e.Email
					break
				}
			}
		}
	}
	if ui.Email == "" {
		ui.Email = fmt.Sprintf("%s@github.user", ui.Login)
	}
	return &ui, nil
}

func fetchGitHubEmails(ctx context.Context, token string) ([]githubEmailInfo, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", "https://api.github.com/user/emails", nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Accept", "application/json")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	var emails []githubEmailInfo
	if err := json.Unmarshal(body, &emails); err != nil {
		return nil, err
	}
	return emails, nil
}

func (a *Auth) OAuthLoginURL(provider, role string, cfg OAuthConfig) string {
	state := fmt.Sprintf("%s:%s", provider, role)
	switch provider {
	case "google":
		return googleRedirectURL(cfg, state)
	case "github":
		return githubRedirectURL(cfg, state)
	default:
		return ""
	}
}

func (a *Auth) HandleOAuthCallback(ctx context.Context, provider, code, state string, cfg OAuthConfig) (*TokenPair, error) {
	parts := strings.SplitN(state, ":", 2)
	role := "consumer"
	if len(parts) == 2 {
		role = parts[1]
	}
	var subject, email, name string
	switch provider {
	case "google":
		ui, err := exchangeGoogleCode(ctx, cfg, code)
		if err != nil {
			return nil, err
		}
		subject = ui.Sub
		email = ui.Email
		name = ui.Name
	case "github":
		ui, err := exchangeGitHubCode(ctx, cfg, code)
		if err != nil {
			return nil, err
		}
		subject = fmt.Sprintf("%d", ui.ID)
		email = ui.Email
		if ui.Name != "" {
			name = ui.Name
		} else {
			name = ui.Login
		}
	default:
		return nil, fmt.Errorf("unsupported provider: %s", provider)
	}
	return a.SocialLogin(ctx, SocialLoginInput{
		Provider: provider,
		Subject:  subject,
		Email:    email,
		Name:    name,
		Role:    role,
	})
}
