package api

import (
	"encoding/base64"
	"encoding/json"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
	"github.com/ssksameer56/CloudIndexer/config"
	"github.com/ssksameer56/CloudIndexer/models"
	"github.com/ssksameer56/CloudIndexer/utils"
)

func processAuthorizationToken(c *gin.Context) {
	authorizationToken := c.Param("code")

	client := utils.HttpClient{
		BaseURL: "https://api.dropbox.com/oauth2/token",
		Client:  &http.Client{},
		Timeout: time.Second * 30,
		Headers: make(map[string]string),
	}

	body := models.DropboxOAuth2Request{
		AuthorizationCode: authorizationToken,
		GrantType:         "authorization_code",
	}
	headers := make(map[string]string)

	auth := config.Config.DropBoxAppKey + ":" + config.Config.DropBoxAppSecret
	basicAuth := base64.StdEncoding.EncodeToString([]byte(auth))
	headers["Authorization"] = "Basic " + basicAuth
	data, err := client.Post("", body, headers, nil)
	if err != nil {
		log.Panic().Err(err).Msgf("cant get access token")
	}
	var response models.DropBoxOAuth2TokenResponse
	err = json.Unmarshal(data, &response)
	if err != nil {
		log.Err(err).Msgf("cant unmarshal access token")
	}
	config.Config.AccessToken = response.AccessToken
}
