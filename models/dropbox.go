package models

type DropBoxFileListRequest struct {
	IncludeDeleted                  bool   `json:"include_deleted,omitempty"`
	IncludeHasExplicitSharedMembers bool   `json:"include_has_explicit_shared_members,omitempty"`
	IncludeMediaInfo                bool   `json:"include_media_info,omitempty"`
	IncludeMountedFolders           bool   `json:"include_mounted_folders,omitempty"`
	IncludeNonDownloadableFiles     bool   `json:"include_non_downloadable_files,omitempty"`
	Path                            string `json:"path,omitempty"`
	Recursive                       bool   `json:"recursive,omitempty"`
}

type DropBoxPollRequest struct {
	Cursor string `json:"cursor,omitempty"`
}

type DropBoxDownloadRequest struct {
	Path string `json:"path,omitempty"`
}

type DropboxCursorRequest struct {
	IncludeDeleted                  bool   `json:"include_deleted,omitempty"`
	IncludeHasExplicitSharedMembers bool   `json:"include_has_explicit_shared_members,omitempty"`
	IncludeMediaInfo                bool   `json:"include_media_info,omitempty"`
	IncludeMountedFolders           bool   `json:"include_mounted_folders,omitempty"`
	IncludeNonDownloadableFiles     bool   `json:"include_non_downloadable_files,omitempty"`
	Path                            string `json:"path,omitempty"`
	Recursive                       bool   `json:"recursive,omitempty"`
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

type DropBoxCursorResponse struct {
	Cursor string `json:"cursor,omitempty"`
}
