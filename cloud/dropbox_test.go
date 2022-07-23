package cloud

import (
	"context"
	"testing"
	"time"

	"github.com/rs/zerolog/log"
	"github.com/ssksameer56/CloudIndexer/config"
	"github.com/stretchr/testify/require"
)

func TestGetFiles(t *testing.T) {
	err := config.LoadConfig()
	if err != nil {
		log.Err(err).Msg("cant load config")
		t.FailNow()
	}
	dropbox := DropBox{
		AuthKey: config.Config.DropboxKey,
		Timeout: time.Minute,
	}
	err = dropbox.Connect(context.Background())
	if err != nil {
		log.Err(err).Msg("cant load dropbox client")
		t.FailNow()
	}
	files, err := dropbox.GetFiles(context.Background(), "/CloudIndexer")
	require.NoError(t, err)
	require.NotEmpty(t, files)
}
