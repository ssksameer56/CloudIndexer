package models

type FileData struct {
	Name string
	Path string
}

type ElasticSearchData struct {
	Name     string `json:"name,omitempty"`
	FilePath string `json:"file_path,omitempty"`
	Text     string `json:"text,omitempty"`
}
