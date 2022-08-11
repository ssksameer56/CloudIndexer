package models

type APISearchResponse struct {
	Name string `json:"name,omitempty"`
	URL  string `json:"url,omitempty"`
}

type APISearchRequest struct {
	Keyword string `form:"keyword,omitempty"`
}

type DropboxOAuth2RefreshRequest struct {
	GrantType    string `json:"grant_type,omitempty"`
	RefreshToken string `json:"refresh_token,omitempty"`
}

type DropboxOAuth2FreshRequest struct {
	AuthorizationCode string `json:"code,omitempty"`
	GrantType         string `json:"grant_type,omitempty"`
	ClientID          string `json:"client_id,omitempty"`
	ClientSecret      string `json:"client_secret,omitempty"`
	RefreshToken      string `json:"refresh_token,omitempty"`
	RedirectURI       string `json:"redirect_uri,omitempty"`
}

type DropBoxOAuth2TokenResponse struct {
	AccessToken string `json:"access_token"`
	ExpiresIn   int64  `json:"expires_in"`
	TokenType   string `json:"token_type"`
	Scope       string `json:"scope"`
	AccountID   string `json:"account_id"`
	UID         string `json:"uid"`
}
