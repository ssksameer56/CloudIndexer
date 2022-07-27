package models

type APISearchResponse struct {
	Name string
	URL  string
}

type DropBoxFileListResponse struct {
	Entries []DropBoxFileMetadata `json:"entries,omitempty"`
	Cursor  string                `json:"cursor,omitempty"`
	HasMore bool                  `json:"has_more,omitempty"`
}

type DropBoxFileMetadata struct {
	Tag         string `json:".tag,omitempty"`
	Name        string `json:"name,omitempty"`
	PathLower   string `json:"path_lower,omitempty"`
	PathDisplay string `json:"path_display,omitempty"`
	ID          string `json:"id,omitempty"`
}

type DropBoxPollResponse struct {
	Changes bool `json:"changes,omitempty"`
}
