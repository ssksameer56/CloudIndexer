package workers

import (
	"context"
	"strings"
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
	changeNotificationChannel  chan models.FolderChangeNotification
	currentPositions           map[string]string
	IndexerNotificationChannel chan models.CloudWatcherNotification
	isWaiting                  map[string]bool
}

func (cw *CloudWatcher) Init(ctx context.Context) error {
	cw.context = ctx
	err := cw.CloudProvider.Connect(cw.context)
	if err != nil {
		log.Err(err).Msg("couldnt start cloud watcher")
		return err
	}
	cw.changeNotificationChannel = make(chan models.FolderChangeNotification, config.Config.BufferSize)
	foldersToWatch := strings.Split(config.Config.Folders, ",")
	cw.currentPositions = make(map[string]string)
	cw.isWaiting = make(map[string]bool)
	for _, folder := range foldersToWatch {
		cursor, err := cw.CloudProvider.GetPointerToPath(cw.context, folder)
		if err != nil {
			log.Err(err).Msg("couldnt get pointer to folder to watch")
			return err
		}
		cw.currentPositions[folder] = cursor
		cw.isWaiting[folder] = false
	}
	err = cw.CloudProvider.Ping(cw.context)
	if err != nil {
		log.Err(err).Msg("couldnt ping dropbox")
		return err
	}
	log.Info().Msg("started cloud worker")
	return nil
}
func (cw *CloudWatcher) Run(wg *sync.WaitGroup) {
	defer wg.Done()

	wnotifs := sync.WaitGroup{}
	for name := range cw.currentPositions {
		wnotifs.Add(1)
		go cw.WaitForNotifcation(&wnotifs, name)
	}

	ticker := time.NewTicker(time.Minute * 10)
	select {
	case <-cw.context.Done():
		log.Info().Str("component", "CloudWatcher").Msg("context done received. exiting cloud watcher loop")
		close(cw.IndexerNotificationChannel)
		return
	case <-ticker.C:
		log.Info().Str("component", "CloudWatcher").Msg("pinging. cloud watcher alive")
	case changeNotif := <-cw.changeNotificationChannel:
		if changeNotif.Change {
			fileList, cursor, err := cw.CloudProvider.GetListofFiles(cw.context, changeNotif.Folder)
			if err != nil {
				log.Err(err).Str("component", "CloudWatcher").Msg("couldnt get latest files")
			}
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
			notif := models.CloudWatcherNotification{
				Folder: changeNotif.Folder,
				Cursor: cursor,
				Data:   newData,
			}
			cw.currentPositions[notif.Folder] = cursor
			cw.isWaiting[notif.Folder] = false
			cw.IndexerNotificationChannel <- notif
		}
	}
	wnotifs.Wait()
}

func (cw *CloudWatcher) WaitForNotifcation(wg *sync.WaitGroup, folder string) {
	defer wg.Done()
	for {
		select {
		case <-cw.context.Done():
			log.Info().Str("component", "CloudWatcher").Msg("wait for notification ended")
			close(cw.changeNotificationChannel)
			return
		default:
			if !cw.isWaiting[folder] {
				go func() {
					cw.CloudProvider.CheckForChange(cw.context, cw.currentPositions[folder], time.Minute*15,
						cw.changeNotificationChannel, folder)
				}()
				cw.isWaiting[folder] = true
			}
		}
	}

}

func (cw *CloudWatcher) Stop() error {
	return nil
}
