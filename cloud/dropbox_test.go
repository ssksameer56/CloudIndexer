package cloud

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/rs/zerolog/log"
	"github.com/ssksameer56/CloudIndexer/config"
	"github.com/ssksameer56/CloudIndexer/models"
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
	files, cursor, err := dropbox.GetListofFiles(context.Background(), "/CloudIndexer")
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

func TestPoll(t *testing.T) {
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
	nChan := make(chan models.FolderChangeNotification, 10)
	cursor := "AAG2DWsHsms6Pl4kNgNowovdWJdTGY0-lSdyUb7fxRdECZCUXscwmy4EdSAbU8WEUs6PQ-YeNF25lKHWJ7TPaQKv3zoUhEAY67OMd1gYSPb7FDqmFW1Z9d0qyXdQuYobYchAiSRVHKTR1QWheLirxup6ykGnAs1alfJTaqK3MsbGrcS0A9Qw5bmHRKLGzj7vgnJxl5pcBk77HEoIQTp-MUK8iIFw-r6QNIFdx2xcfYF0Tg"
	dropbox.CheckForChange(context.Background(), cursor, time.Minute, nChan, "/CloudIndexer")
	for c := range nChan {
		fmt.Printf("%s %v %v", time.Now().String(), c.Folder, c.Change)
	}
}
