package handlers

import (
	"context"

	"github.com/rs/zerolog/log"
	"github.com/ssksameer56/CloudIndexer/cloud"
	"github.com/ssksameer56/CloudIndexer/elasticservice"
	"github.com/ssksameer56/CloudIndexer/models"
)

type SearchHandler struct {
	CloudProvider   cloud.Cloud
	ESSearchService elasticservice.ElasticSearchService
}

func (sh *SearchHandler) SearchText(ctx context.Context, index, keyword string) ([]models.TextStoreModel, error) {
	cctx, cancel := context.WithCancel(ctx)
	defer cancel()
	results, err := sh.ESSearchService.Search(cctx, index, keyword)
	if err != nil {
		log.Err(err).Msgf("error when trying to search ES for %s %s", index, keyword)
		return []models.TextStoreModel{}, err
	}
	var files []models.TextStoreModel
	for _, res := range results.Hits.Hits {
		file := models.TextStoreModel{}
		file.FilePath = res.Source.FilePath
		file.Name = res.Source.Name
		files = append(files, file)
	}
	return files, nil
}
