package utils

import (
	"io"

	"github.com/rs/zerolog/log"
)

func ReadFile(reader io.ReadCloser) ([]byte, error) {
	/*
	   You convert a Reader to bytes, by reading it. There's not really a more efficient way to do it.

	   body, err := ioutil.ReadAll(r.Body)
	   If you are unconditionally transferring bytes from an io.Reader to an io.Writer, you can just use io.Copy
	*/
	body, err := io.ReadAll(reader)
	if err != nil {
		log.Err(err).Msg("couldnt read data from reader")
	}

	log.Info().Msgf("read a file %s with size %d", len(body))
	return body, nil
}
