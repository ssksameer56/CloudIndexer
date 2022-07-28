package handlers

import (
	"context"

	"github.com/rs/zerolog/log"
	"github.com/ssksameer56/CloudIndexer/cloud"
	"github.com/ssksameer56/CloudIndexer/elasticservice"
	"github.com/ssksameer56/CloudIndexer/models"
)

type IndexHandler struct {
	Service         cloud.Cloud
	ESSearchService elasticservice.ElasticSearchService
	Index           string
}

func (ih *IndexHandler) IndexStuff(ctx context.Context, data models.TextStoreModel) (string, error) {
	cctx, cancel := context.WithCancel(ctx)
	defer cancel()
	res, err := ih.ESSearchService.Index(cctx, ih.Index, data)
	if err != nil {
		log.Err(err).Msgf("couldnt index file %s %s", data.FilePath, data.Name)
		return "", nil
	}
	id := res.ID
	return id, nil
}
