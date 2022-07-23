package utils

import (
	"errors"
	"io/ioutil"
	"net/http"
	"net/url"
	"time"

	"github.com/rs/zerolog/log"
)

type HttpClient struct {
	Client  *http.Client
	Timeout time.Duration
}

func (hc *HttpClient) Get(reqURL string, headers map[string]interface{}) ([]byte, error) {
	var response *http.Response
	var request http.Request
	var err error

	request.URL, err = url.Parse(reqURL)
	if err != nil {
		log.Error().Err(err).Msgf("cant parse URL %s", reqURL)
	}
	for name, value := range headers {
		request.Header.Add(name, value.(string))
	}

	hc.Client.Timeout = hc.Timeout
	response, err = hc.Client.Do(&request)
	if err != nil || response.StatusCode != 200 {
		log.Error().Err(err).Msgf("failed to send request")
		if response.StatusCode != 200 {
			return []byte{}, errors.New("no 200 response")
		} else {
			return []byte{}, err
		}
	}
	defer response.Body.Close()
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}
	return body, nil
}
