package models

type AppConfig struct {
	DropboxKey       string `json:"DROPBOX_API_KEY,omitempty"`
	ElasticSearchURL string `json:"ELASTICSEARCH_URL,omitempty"`
	BufferSize       int    `json:"BUFFER_SIZE,omitempty"`
}
