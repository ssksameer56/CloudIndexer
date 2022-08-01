package cloud

import (
	"context"
	"time"

	"github.com/ssksameer56/CloudIndexer/models"
)

type Cloud interface {
	GetListofFiles(ctx context.Context, name string) ([]models.FileData, string, error)
	Connect(ctx context.Context) error
	CheckForChange(ctx context.Context, cursor string, timeout time.Duration,
		notifcationChannel chan models.FolderChangeNotification, folder string)
	GetPointerToPath(ctx context.Context, path string) (string, error)
	DownloadFile(ctx context.Context, filePath string) ([]byte, error)
}
