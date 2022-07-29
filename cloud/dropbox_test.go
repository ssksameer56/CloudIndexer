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
	files, cursor, err := dropbox.GetFiles(context.Background(), "/CloudIndexer")
	require.NoError(t, err)
	require.NotEmpty(t, files)
	require.NotEmpty(t, cursor)
}

func TestDownload(t *testing.T) {
	err := config.LoadConfig()
	if err != nil {
		log.Err(err).Str("component", "Dropbox Test").Msg("cant load config")
		t.FailNow()
	}
	dropbox := DropBox{
		AuthKey: config.Config.DropboxKey,
		Timeout: time.Minute,
	}
	err = dropbox.Connect(context.Background())
	if err != nil {
		log.Err(err).Str("component", "Dropbox Test").Msg("cant load dropbox client")
		t.FailNow()
	}
	files, err := dropbox.DownloadFile(context.Background(), "/CloudIndexer/hello.txt")
	require.NoError(t, err)
	require.NotEmpty(t, files)
}

func TestCursor(t *testing.T) {
	err := config.LoadConfig()
	if err != nil {
		log.Err(err).Str("component", "Dropbox Test").Msg("cant load config")
		t.FailNow()
	}
	dropbox := DropBox{
		AuthKey: config.Config.DropboxKey,
		Timeout: time.Minute,
	}
	err = dropbox.Connect(context.Background())
	if err != nil {
		log.Err(err).Str("component", "Dropbox Test").Msg("cant load dropbox client")
		t.FailNow()
	}
	files, err := dropbox.GetPointerToPath(context.Background(), "/CloudIndexer")
	require.NoError(t, err)
	require.NotEmpty(t, files)
}
