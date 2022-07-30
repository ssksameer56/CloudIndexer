package workers

import (
	"context"
	"sync"
	"time"

	"github.com/rs/zerolog/log"
	"github.com/ssksameer56/CloudIndexer/cloud"
	"github.com/ssksameer56/CloudIndexer/config"
	"github.com/ssksameer56/CloudIndexer/models"
)

type CloudWatcher struct {
	CloudProvider              cloud.Cloud
	context                    context.Context
	changeNotificationChannel  chan bool
	FolderToWatch              string
	currentPosition            string
	IndexerNotificationChannel chan models.CloudWatcherNotification
}

func (cw *CloudWatcher) Init(ctx context.Context) error {
	cw.context = ctx
	err := cw.CloudProvider.Connect(cw.context)
	if err != nil {
		log.Err(err).Msg("couldnt start cloud watcher")
		return err
	}
	cw.changeNotificationChannel = make(chan bool, config.Config.BufferSize)
	cursor, err := cw.CloudProvider.GetPointerToPath(cw.context, cw.FolderToWatch)
	if err != nil {
		log.Err(err).Msg("couldnt get pointer to folder to watch")
		return err
	}
	cw.currentPosition = cursor
	return nil
}
func (cw *CloudWatcher) Run(wg *sync.WaitGroup) {
	go cw.WaitForNotifcation()
	defer wg.Done()
	ticker := time.NewTicker(time.Minute * 10)
	select {
	case <-cw.context.Done():
		log.Info().Str("component", "CloudWatcher").Msg("context done received. exiting cloud watcher loop")
		close(cw.IndexerNotificationChannel)
	case <-ticker.C:
		log.Info().Str("component", "CloudWatcher").Msg("pinging. cloud watcher alive")
	case <-cw.changeNotificationChannel:
		//TODO: update stuff to download files while call is made for folder
		fileList, cursor, err := cw.CloudProvider.GetListofFiles(cw.context, cw.FolderToWatch)
		newData := make([]models.TextStoreModel, len(fileList))
		for i, file := range fileList {
			go func(i int, file models.FileData) {
				data, err := cw.CloudProvider.DownloadFile(cw.context, file.Path)
				if err != nil {
					log.Err(err).Msgf("error when downloading %s", file.Path)
				}
				newData[i] = models.TextStoreModel{
					Name:     file.Name,
					FilePath: file.Path,
					Text:     string(data),
				}
			}(i, file)
		}
		if err != nil {
			log.Err(err).Str("component", "CloudWatcher").Msg("couldnt get latest cursor")
		}
		notif := models.CloudWatcherNotification{
			Folder: cw.FolderToWatch,
			Cursor: cursor,
			Data:   newData,
		}
		cw.currentPosition = cursor
		cw.IndexerNotificationChannel <- notif
	}

}

func (cw *CloudWatcher) WaitForNotifcation() {
	select {
	case <-cw.context.Done():
		log.Info().Str("component", "CloudWatcher").Msg("wait for notification ended")
		close(cw.changeNotificationChannel)
		return
	default:
		cw.CloudProvider.CheckForChange(cw.context, cw.currentPosition, time.Hour, cw.changeNotificationChannel)
	}
}

func (cw *CloudWatcher) Stop() error {
	return nil
}
