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
