package config

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/rs/zerolog/log"
	"github.com/ssksameer56/CloudIndexer/models"
	"github.com/ssksameer56/CloudIndexer/utils"
)

var Config models.AppConfig

func LoadConfig() error {
	raw, err := ioutil.ReadFile("../config/config.json")
	if err != nil {
		log.Panic().Msg("Error occured while reading config")
		return err
	}
	json.Unmarshal(raw, &Config)
	OpenOAuth2TokenPopup()
	return nil
}

func GetAccessToken() (string, error) {
	if Config.AccessToken == "" {
		OpenOAuth2TokenPopup()
		return Config.AccessToken, errors.New(models.NoAccessToken)
	} else {
		return Config.AccessToken, nil
	}
}

func OpenOAuth2TokenPopup() {
	url := "https://www.dropbox.com/oauth2/authorize?client_id=%s&redirect_uri=%s&response_type=code"
	url = fmt.Sprintf(url, Config.DropBoxAppKey, Config.RedirectURI)
	client := utils.HttpClient{
		BaseURL: url,
		Client:  &http.Client{},
		Timeout: time.Second * 30,
		Headers: make(map[string]string),
	}
	_, err := client.Get("")
	if err != nil {
		log.Err(err).Msg("couldnt open login popup")
	}
}
