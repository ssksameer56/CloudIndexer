package models

type AppConfig struct {
	DropboxKey         string `json:"DROPBOX_API_KEY,omitempty"`
	ElasticSearchURL   string `json:"ELASTICSEARCH_URL,omitempty"`
	BufferSize         int    `json:"BUFFER_SIZE,omitempty"`
	Folders            string `json:"FOLDERS_TO_WATCH"`
	DropBoxAppKey      string `json:"DROPBOX_APP_KEY"`
	DropBoxAppSecret   string `json:"DROPBOX_APP_SECRET"`
	RedirectURI        string `json:"REDIRECT_URI"`
	AccessToken        string
	AuthorizationToken string `json:"DROPBOX_AUTHORIZATION_CODE"`
	RefreshToken       string `json:"DROPBOX_REFRESH_TOKEN"`
}
