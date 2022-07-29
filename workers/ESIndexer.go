package workers

import (
	"context"
	"sync"
	"time"

	"github.com/rs/zerolog/log"
	"github.com/ssksameer56/CloudIndexer/elasticservice"
	"github.com/ssksameer56/CloudIndexer/models"
)

type ESWorker struct {
	Service                    elasticservice.ElasticSearchService
	AuthCode                   string
	Context                    context.Context
	IndexerNotificationChannel chan models.CloudWatcherNotification
}

func (cw *ESWorker) Init(ctx context.Context) error {
	cw.Service = elasticservice.ElasticSearchService{}
	err := cw.Service.Connect()
	if err != nil {
		log.Err(err).Str("component", "ElasticSearchIndexer").Msg("couldnt establish connection to ES")
		return err
	}
	return nil
}
func (cw *ESWorker) Run(wg *sync.WaitGroup) {
	defer wg.Done()
	ticker := time.NewTicker(time.Minute * 10)
	select {
	case <-cw.Context.Done():
		log.Info().Str("component", "ElasticSearchIndexer").Msg("context done received. exiting es indexer loop")
	case <-ticker.C:
		log.Info().Str("component", "ElasticSearchIndexer").Msg("pinging. es indexer alive")
	case data := <-cw.IndexerNotificationChannel:

	}

}

func (cw *ESWorker) Stop() error {
	return nil
}
