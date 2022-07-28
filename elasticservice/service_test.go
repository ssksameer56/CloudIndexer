package elasticservice

import (
	"context"
	"testing"

	"github.com/rs/zerolog/log"
	"github.com/ssksameer56/CloudIndexer/config"
	"github.com/ssksameer56/CloudIndexer/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type ElasticSearchTestSuite struct {
	suite.Suite
	Service ElasticSearchService
}

func (ets *ElasticSearchTestSuite) SetupTest() {
	_ = assert.New(ets.T())
	config.LoadConfig()
	ets.Service = ElasticSearchService{}
	err := ets.Service.Connect()
	if err != nil {
		log.Err(err).Msg("couldnt establish connection to ES")
	}
}

func (ets *ElasticSearchTestSuite) TestIndex() {
	data := models.TextStoreModel{
		Name:     "test file",
		FilePath: "some/weird/stuff/",
		Text:     "Quick brown fox jumps over the lazy dog",
	}
	res, err := ets.Service.Index(context.Background(), "cloud-indexer", data)
	require.NoError(ets.T(), err, "error while indexing data to es")
	require.NotEmpty(ets.T(), res)
}

func (ets *ElasticSearchTestSuite) TestSearch() {
	res, err := ets.Service.Search(context.Background(), "cloud-indexer", "brown")
	require.NoError(ets.T(), err, "error while indexing data to es")
	require.NotEmpty(ets.T(), res)
}

func TestElasticSearchTestSuite(t *testing.T) {
	suite.Run(t, new(ElasticSearchTestSuite))
}
