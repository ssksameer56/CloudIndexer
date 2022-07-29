package models

type APISearchResponse struct {
	Name string `json:"name,omitempty"`
	URL  string `json:"url,omitempty"`
}

type APISearchRequest struct {
	Keyword string `json:"keyword,omitempty"`
}
