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

func (db *DropBox) GetFiles(ctx context.Context, path string) ([]models.FileData, error) {
	_, cancel := context.WithCancel(ctx)
	defer cancel()
	url := "files/list_folder"

	body := models.DropBoxFileListRequest{
		Path: path,
	}
	data, err := db.client.Post(url, body)
	if err != nil {
		log.Error().Err(err).Msgf("cant get files for %s", path)
		return []models.FileData{}, err
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
	return fileNames, err
}

func (db *DropBox) Connect(ctx context.Context) error {
	db.AuthKey = config.Config.DropboxKey
	db.client = utils.HttpClient{
		BaseURL: "https://api.dropboxapi.com/2/",
		Client:  &http.Client{},
		Timeout: time.Second * 30,
		Headers: make(map[string]string),
	}
	db.client.Headers["Authorization"] = "Bearer " + db.AuthKey
	return nil
}
