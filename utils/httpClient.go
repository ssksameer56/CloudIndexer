package utils

import (
	"bytes"
	"encoding/json"
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
	BaseURL string
}

func (hc *HttpClient) Get(reqURL string, headers map[string]string) ([]byte, error) {
	var response *http.Response
	var request http.Request
	var err error

	reqURL = hc.BaseURL + reqURL
	request.URL, err = url.Parse(reqURL)
	if err != nil {
		log.Error().Err(err).Msgf("cant parse URL %s", reqURL)
	}
	for name, value := range headers {
		request.Header.Add(name, value)
	}
	request.Header.Add("Content-Type", "application/json")

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

func (hc *HttpClient) Post(reqURL string, headers map[string]string, body interface{}) ([]byte, error) {
	var response *http.Response
	var request *http.Request
	var err error

	if err != nil {
		log.Error().Err(err).Msgf("cant parse URL %s", reqURL)
	}

	reqBody, err := json.Marshal(body)
	if err != nil {
		log.Err(err).Msgf("cannot marshal body to do post request %s", reqURL)
		return []byte{}, err
	}

	url := hc.BaseURL + reqURL
	request, err = http.NewRequest("POST", url, bytes.NewBuffer(reqBody))
	if err != nil {
		log.Err(err).Msgf("cannot create post request %s", reqURL)
		return nil, err
	}
	for name, value := range headers {
		request.Header.Add(name, value)
	}
	request.Header.Add("Content-Type", "application/json")

	response, err = hc.Client.Do(request)

	if err != nil || response.StatusCode != 200 {
		log.Error().Err(err).Msgf("failed to send request")
		if response.StatusCode != 200 {
			return []byte{}, errors.New(response.Status)
		} else {
			return []byte{}, err
		}
	}
	defer response.Body.Close()
	body, err = ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}
	return body.([]byte), nil
}
