package main

import (
	"github.com/rs/zerolog/log"
	"github.com/ssksameer56/CloudIndexer/config"
)

func main() {
	err := config.LoadConfig()
	if err != nil {
		log.Panic().Err(err).Msg("Cant load config")
	}

}
