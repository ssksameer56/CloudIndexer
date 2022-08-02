package workers

import (
	"context"
	"errors"
	"sync"
	"time"

	"github.com/rs/zerolog/log"
	"github.com/ssksameer56/CloudIndexer/elasticservice"
	"github.com/ssksameer56/CloudIndexer/models"
)

type ESWorker struct {
	Service                    elasticservice.ElasticSearchService
	context                    context.Context
	IndexerNotificationChannel chan models.CloudWatcherNotification
}

func (esw *ESWorker) Init(ctx context.Context) error {
	esw.Service = elasticservice.ElasticSearchService{}
	err := esw.Service.Connect()
	if err != nil {
		log.Err(err).Str("component", "ElasticSearchIndexer").Msg("couldnt establish connection to ES")
		return err
	}
	return nil
}
func (esw *ESWorker) Run(wg *sync.WaitGroup) {
	defer wg.Done()
	ticker := time.NewTicker(time.Minute * 10)
	select {
	case <-esw.context.Done():
		log.Info().Str("component", "ElasticSearchIndexer").Msg("context done received. exiting es indexer loop")
		return
	case <-ticker.C:
		log.Info().Str("component", "ElasticSearchIndexer").Msg("pinging. es indexer alive")
	case data := <-esw.IndexerNotificationChannel:
		for _, item := range data.Data {
			go func(item models.TextStoreModel) {
				res, err := esw.Service.Index(esw.context, data.Folder, item)
				if err != nil {
					log.Err(err).Str("component", "ElasticSearchIndexer").Msgf("couldnt index data %s", item.FilePath)
				}
				if res.Result != models.ESIndexCreated || res.Result != models.ESIndexUpdated {
					log.Err(errors.New("couldnt created")).Str("component", "ElasticSearchIndexer").
						Msgf("couldnt index data %s", item.FilePath)
				}
			}(item)
		}
	}

}

func (esw *ESWorker) Stop() error {
	log.Info().Str("component", "ElasticSearchIndexer").
		Msgf("exiting es indexer loop")
	return nil
}
