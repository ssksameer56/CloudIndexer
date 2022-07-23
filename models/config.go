package models

type AppConfig struct {
	DropboxKey string `json:"DROPBOX_API_KEY,omitempty"`
	UniDocKey  string `json:"UNIDOC_API_KEY,omitempty"`
}
