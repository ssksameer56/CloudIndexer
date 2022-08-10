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
	positionMutex              sync.RWMutex
	waitMutex                  sync.RWMutex
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
	for {
		select {
		case <-cw.context.Done():
			log.Info().Str("component", "CloudWatcher").Msg("context done received. exiting cloud watcher loop")
			close(cw.IndexerNotificationChannel)
			wnotifs.Wait()
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
				wg := sync.WaitGroup{}
				wg.Add(len(fileList))
				for i, file := range fileList {
					go func(wg *sync.WaitGroup, i int, file models.FileData, dataArray *[]models.TextStoreModel) {
						defer wg.Done()
						data, err := cw.CloudProvider.DownloadFile(cw.context, file.Path)
						if err != nil {
							log.Err(err).Str("component", "CloudWatcher").Msgf("error when downloading %s", file.Path)
						}
						(*dataArray)[i] = models.TextStoreModel{
							Name:     file.Name,
							FilePath: file.Path,
							Text:     string(data),
						}
					}(&wg, i, file, &newData)
				}
				wg.Wait()
				notif := models.CloudWatcherNotification{
					Folder: changeNotif.Folder,
					Cursor: cursor,
					Data:   newData,
				}
				cw.waitMutex.Lock()
				cw.positionMutex.Lock()
				cw.currentPositions[notif.Folder] = cursor
				cw.isWaiting[notif.Folder] = false
				cw.waitMutex.Unlock()
				cw.positionMutex.Unlock()
				cw.IndexerNotificationChannel <- notif
				log.Info().Msgf("NOTIF %v", changeNotif)
			}
		}
	}
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
			cw.waitMutex.RLock()
			flag := cw.isWaiting[folder]
			cw.waitMutex.RUnlock()
			if !flag {
				go func() {
					cw.CloudProvider.CheckForChange(cw.context, cw.currentPositions[folder], time.Minute*15,
						cw.changeNotificationChannel, folder)
				}()
				cw.waitMutex.Lock()
				cw.isWaiting[folder] = true
				cw.waitMutex.Unlock()
			}
		}
	}
}

func (cw *CloudWatcher) Stop() error {
	return nil
}
