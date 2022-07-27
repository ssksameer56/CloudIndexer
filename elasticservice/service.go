package elasticservice

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/elastic/go-elasticsearch/v8"
	"github.com/elastic/go-elasticsearch/v8/esapi"
	"github.com/elastic/go-elasticsearch/v8/esutil"
	"github.com/rs/zerolog/log"
	"github.com/ssksameer56/CloudIndexer/config"
	"github.com/ssksameer56/CloudIndexer/models"
)

type ElasticSearchService struct {
	Conn *elasticsearch.Client
}

func (es *ElasticSearchService) Connect() error {
	var err error
	es.Conn, err = elasticsearch.NewClient(elasticsearch.Config{
		Addresses: []string{config.Config.ElasticSearchURL},
	})
	if err != nil {
		log.Panic().Err(err).Msg("couldnt connect to ES")
		return err
	}
	return nil
}

func (es *ElasticSearchService) Search(ctx context.Context, index, keyword string) ([]models.ESSearchResponse, error) {
	cctx, cancel := context.WithCancel(ctx)
	defer cancel()
	req := models.SearchRequest{
		Query: models.Query{
			Match: models.Match{
				Text: keyword,
			},
		},
	}
	body := esutil.NewJSONReader(&req)
	res, err := es.Conn.Search(
		es.Conn.Search.WithContext(cctx),
		es.Conn.Search.WithIndex(index),
		es.Conn.Search.WithBody(body),
	)
	if err != nil {
		log.Err(err).Msg("couldnt search to ES")
		return []models.ESSearchResponse{}, err
	}
	defer res.Body.Close()
	resBody, err := ioutil.ReadAll(res.Body)

	if res.StatusCode != http.StatusOK {
		log.Err(err).Msg("couldnt get success response from ES")
		return []models.ESSearchResponse{}, err
	}
	var ESResults []models.ESSearchResponse
	if err != nil {
		log.Err(err).Msg("couldnt read response body")
		return []models.ESSearchResponse{}, err
	}
	err = json.Unmarshal(resBody, &ESResults)
	if err != nil {
		log.Err(err).Msg("couldnt unmarshal res body")
		return []models.ESSearchResponse{}, err
	}
	return ESResults, nil
}

func (es *ElasticSearchService) Index(ctx context.Context, index string, data models.TextStoreModel) (models.ESIndexResponse, error) {
	cctx, cancel := context.WithCancel(ctx)
	defer cancel()

	dataJSON := esutil.NewJSONReader(data)
	req := esapi.IndexRequest{
		Index:   index,
		Body:    dataJSON,
		Refresh: "true",
	}
	res, err := req.Do(cctx, es.Conn)
	if err != nil {
		log.Err(err).Msg("couldnt index to ES")
		return models.ESIndexResponse{}, err
	}
	defer res.Body.Close()
	resBody, err := ioutil.ReadAll(res.Body)

	if res.StatusCode != http.StatusOK {
		log.Err(err).Msg("couldnt get success response from ES")
		return models.ESIndexResponse{}, err
	}
	var ESResults models.ESIndexResponse
	if err != nil {
		log.Err(err).Msg("couldnt read response body")
		return models.ESIndexResponse{}, err
	}
	err = json.Unmarshal(resBody, &ESResults)
	if err != nil {
		log.Err(err).Msg("couldnt unmarshal res body")
		return models.ESIndexResponse{}, err
	}
	return ESResults, nil
}
