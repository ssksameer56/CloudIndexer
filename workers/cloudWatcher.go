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
	AuthCode                   string
	Context                    context.Context
	ChangeNotificationChannel  chan bool
	FolderToWatch              string
	CurrentPosition            string
	IndexerNotificationChannel chan models.CloudWatcherNotification
}

func (cw *CloudWatcher) Init(ctx context.Context) error {
	cw.Context = ctx
	err := cw.CloudProvider.Connect(cw.Context)
	if err != nil {
		log.Err(err).Msg("couldnt start cloud watcher")
		return err
	}
	cw.ChangeNotificationChannel = make(chan bool, config.Config.BufferSize)
	cw.IndexerNotificationChannel = make(chan models.CloudWatcherNotification, config.Config.BufferSize)
	return nil
}
func (cw *CloudWatcher) Run(wg *sync.WaitGroup) {
	go cw.WaitForNotifcation()
	defer wg.Done()
	ticker := time.NewTicker(time.Minute * 10)
	select {
	case <-cw.Context.Done():
		log.Info().Str("component", "CloudWatcher").Msg("context done received. exiting cloud watcher loop")
	case <-ticker.C:
		log.Info().Str("component", "CloudWatcher").Msg("pinging. cloud watcher alive")
	case <-cw.ChangeNotificationChannel:
		cursor, err := cw.CloudProvider.GetPointerToPath(cw.Context, cw.FolderToWatch)
		if err != nil {
			log.Err(err).Str("component", "CloudWatcher").Msg("couldnt get latest cursor")
		}
		notif := models.CloudWatcherNotification{
			Folder: cw.FolderToWatch,
			Cursor: cursor,
		}
		cw.CurrentPosition = cursor
		cw.IndexerNotificationChannel <- notif
	}

}

func (cw *CloudWatcher) WaitForNotifcation() {
	select {
	case <-cw.Context.Done():
		log.Info().Str("component", "CloudWatcher").Msg("wait for notification ended")
		return
	default:
		cw.CloudProvider.CheckForChange(cw.Context, cw.CurrentPosition, time.Hour, cw.ChangeNotificationChannel)
	}
}

func (cw *CloudWatcher) Stop() error {
	return nil
}
