package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
	"github.com/ssksameer56/CloudIndexer/handlers"
	"github.com/ssksameer56/CloudIndexer/models"
)

type SearchController struct {
	Handler handlers.SearchHandler
}

func (sc *SearchController) Search(c *gin.Context) {
	var request models.APISearchRequest
	if err := c.BindQuery(&request); err != nil {
		log.Err(err).Str("component", "SearchController").Msgf("cant bind request for search")
		c.JSON(http.StatusBadRequest, gin.H{})
	}
	data, err := sc.Handler.SearchText(c.Request.Context(), "cloud-indexer", request.Keyword)
	if err != nil {
		log.Err(err).Str("component", "SearchController").Msgf("cant bind request for search")
		c.JSON(http.StatusInternalServerError, gin.H{"error": err})
	}
	var responses []models.APISearchResponse
	for _, item := range data {
		responses = append(responses, models.APISearchResponse{
			Name: item.Name,
			URL:  item.FilePath,
		})
	}
	c.JSON(http.StatusOK, responses)
}
