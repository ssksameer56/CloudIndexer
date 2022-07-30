package cloud

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/rs/zerolog/log"
	"github.com/ssksameer56/CloudIndexer/config"
	"github.com/ssksameer56/CloudIndexer/models"
	"github.com/ssksameer56/CloudIndexer/utils"
)

type DropBox struct {
	client  utils.HttpClient
	AuthKey string
	Timeout time.Duration
}

func (db *DropBox) GetListofFiles(ctx context.Context, folderName string) ([]models.FileData, string, error) {
	_, cancel := context.WithCancel(ctx)
	defer cancel()
	url := "https://api.dropboxapi.com/2/files/list_folder"

	body := models.DropBoxFileListRequest{
		Path: folderName,
	}
	headers := make(map[string]string)
	headers["Content-Type"] = "application/json"
	data, err := db.client.Post(url, body, headers, nil)
	if err != nil {
		log.Error().Err(err).Msgf("cant get files for %s", folderName)
		return []models.FileData{}, "", err
	}
	var response models.DropBoxFileListResponse
	err = json.Unmarshal(data, &response)
	fileNames := []models.FileData{}
	for _, item := range response.Entries {
		fileNames = append(fileNames, models.FileData{
			Name: item.Name,
			Path: item.PathLower,
		})
	}
	return fileNames, response.Cursor, err
}

//TODO : test for poll change if its working
//cursors are unique hashes that point to a folder genereated by dropbox
func (db *DropBox) CheckForChange(ctx context.Context, cursor string, timeout time.Duration,
	notifcationChannel chan bool) {
	_, cancel := context.WithCancel(ctx)
	defer cancel()
	url := "https://notify.dropboxapi.com/2/files/list_folder/longpoll"
	body := models.DropBoxPollRequest{
		Cursor: cursor,
	}
	headers := make(map[string]string)
	headers["Content-Type"] = "application/json"
	data, err := db.client.Post(url, body, headers, &timeout)
	if err != nil {
		log.Error().Err(err).Str("component", "Dropbox").Msgf("cant poll for %s", cursor)
		return
	}
	var response models.DropBoxPollResponse
	err = json.Unmarshal(data, &response)
	if err != nil {
		log.Error().Err(err).Str("component", "Dropbox").Msgf("cant unmarshal result from api call for %s", cursor)
		return
	}
	notifcationChannel <- true
}

func (db *DropBox) DownloadFile(ctx context.Context, filePath string) ([]byte, error) {
	_, cancel := context.WithCancel(ctx)
	defer cancel()
	url := "https://content.dropboxapi.com/2/files/download"
	body := models.DropBoxDownloadRequest{
		Path: filePath,
	}

	path, _ := json.Marshal(body)
	headers := make(map[string]string)
	headers["Dropbox-API-Arg"] = string(path)

	data, err := db.client.Post(url, body, headers, nil)
	if err != nil {
		log.Error().Err(err).Str("component", "Dropbox").Msgf("cant download for %s", filePath)
		return []byte{}, err
	}
	return data, nil
}

//Get only cursor
func (db *DropBox) GetPointerToPath(ctx context.Context, path string) (string, error) {
	_, cancel := context.WithCancel(ctx)
	defer cancel()
	url := "https://api.dropboxapi.com/2/files/list_folder/get_latest_cursor"
	body := models.DropboxCursorRequest{
		Path: path,
	}

	headers := make(map[string]string)
	headers["Content-Type"] = "application/json"
	data, err := db.client.Post(url, body, headers, nil)
	if err != nil {
		log.Error().Err(err).Str("component", "Dropbox").Msgf("cant get cursor for %s", path)
		return "", err
	}
	var response models.DropBoxCursorResponse
	err = json.Unmarshal(data, &response)
	if err != nil {
		log.Error().Err(err).Str("component", "Dropbox").Msgf("cant unmarshal result from api call for %s", path)
		return "", nil
	}
	return response.Cursor, nil
}

func (db *DropBox) Connect(ctx context.Context) error {
	db.AuthKey = config.Config.DropboxKey
	db.client = utils.HttpClient{
		BaseURL: "",
		Client:  &http.Client{},
		Timeout: time.Second * 30,
		Headers: make(map[string]string),
	}
	db.client.Headers["Authorization"] = "Bearer " + db.AuthKey
	return nil
}
