package config

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/rs/zerolog/log"
	"github.com/ssksameer56/CloudIndexer/models"
	"github.com/ssksameer56/CloudIndexer/utils"
)

var Config models.AppConfig

func LoadConfig() error {
	path := ""
	if flag.Lookup("test.v") == nil {
		path = "config/config.json"
	} else {
		path = "config.json"
	}
	raw, err := ioutil.ReadFile(path)
	if err != nil {
		log.Panic().Err(err).Msg("Error occured while reading config")
		return err
	}
	json.Unmarshal(raw, &Config)
	getAccessToken()
	return nil
}

func GetAccessToken() (string, error) {
	if Config.AccessToken == "" {
		getAccessToken()
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

func getAccessToken() {
	BaseURL := "https://api.dropbox.com/oauth2/token"
	client := http.Client{}

	body := models.DropboxOAuth2RefreshRequest{
		RefreshToken: Config.RefreshToken,
		GrantType:    "refresh_token",
	}
	data := url.Values{}
	data.Set("refresh_token", body.RefreshToken)
	data.Set("grant_type", body.GrantType)
	dataURL := data.Encode()

	auth := Config.DropBoxAppKey + ":" + Config.DropBoxAppSecret
	basicAuth := base64.StdEncoding.EncodeToString([]byte(auth))
	r, _ := http.NewRequest(http.MethodPost, BaseURL, strings.NewReader(dataURL))
	r.Header.Add("Authorization", "Basic "+basicAuth)
	r.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	resp, err := client.Do(r)
	if err != nil {
		log.Panic().Err(err).Msg("didnt get auth token")
	}
	if resp.StatusCode != http.StatusOK {
		log.Panic().Err(err).Msg("didnt get auth token 200 resp")
	}
	defer resp.Body.Close()
	respData, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Panic().Err(err).Msg("failed to read response")
	}
	var response models.DropBoxOAuth2TokenResponse
	err = json.Unmarshal(respData, &response)
	if err != nil {
		log.Err(err).Msgf("cant unmarshal access token")
	}
	Config.AccessToken = response.AccessToken
}
