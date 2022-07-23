package config

import (
	"encoding/json"
	"io/ioutil"
	"log"

	"github.com/ssksameer56/CloudIndexer/models"
)

var Config models.AppConfig

func LoadConfig() error {
	raw, err := ioutil.ReadFile("../config/config.json")
	if err != nil {
		log.Println("Error occured while reading config")
		return err
	}
	json.Unmarshal(raw, &Config)
	return nil
}
