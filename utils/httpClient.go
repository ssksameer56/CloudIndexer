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
	Headers map[string]string
}

func (hc *HttpClient) Get(reqURL string) ([]byte, error) {
	var response *http.Response
	var request http.Request
	var err error

	reqURL = hc.BaseURL + reqURL
	request.URL, err = url.Parse(reqURL)
	if err != nil {
		log.Error().Err(err).Msgf("cant parse URL %s", reqURL)
	}
	for name, value := range hc.Headers {
		request.Header.Add(name, value)
	}
	request.Header.Add("Content-Type", "application/json")

	hc.Client.Timeout = hc.Timeout
	response, err = hc.Client.Do(&request)
	if err != nil {
		log.Error().Err(err).Msgf("failed to send request")
		return []byte{}, err
	}
	defer response.Body.Close()
	rbody, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Error().Err(err).Msgf("failed to read response")
		return []byte{}, err
	}
	if response.StatusCode != 200 {
		log.Error().Err(err).Msgf("failed to get response %s", rbody)
		return []byte{}, errors.New(response.Status)
	}
	return rbody, nil
}

func (hc *HttpClient) Post(reqURL string, body interface{}, headers map[string]string, timeout *time.Duration) ([]byte, error) {
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
	for name, value := range hc.Headers {
		request.Header.Add(name, value)
	}
	for name, value := range headers {
		request.Header.Add(name, value)
	}
	if timeout != nil {
		hc.Client.Timeout = *timeout
	}
	response, err = hc.Client.Do(request)

	if err != nil {
		log.Error().Err(err).Msgf("failed to send request")
		return []byte{}, err
	}
	defer response.Body.Close()
	rbody, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Error().Err(err).Msgf("failed to read response")
		return []byte{}, err
	}
	if response.StatusCode != 200 {
		log.Error().Msgf("failed to get response %s", rbody)
		return []byte{}, errors.New(response.Status)
	}

	if err != nil {
		return nil, err
	}
	return rbody, nil
}
