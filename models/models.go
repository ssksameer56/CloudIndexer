package models

type FileData struct {
	Name string
	Path string
}

type TextStoreModel struct {
	Name     string `json:"name,omitempty"`
	FilePath string `json:"file_path,omitempty"`
	Text     string `json:"text,omitempty"`
}

type CloudWatcherNotification struct {
	Cursor string           `json:"cursor,omitempty"`
	Folder string           `json:"file_path,omitempty"` //TODO: use Folder for indexing
	Data   []TextStoreModel `json:"data,omitempty"`
}
