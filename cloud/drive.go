package cloud

import (
	"context"

	"google.golang.org/api/drive/v3"
)

type GoogleDrive struct {
	DriveService drive.Drive
}

func (gd *GoogleDrive) GetFiles(ctx context.Context, name string) ([]string, error) {
	return []string{}, nil
}
func (gd *GoogleDrive) Connect(ctx context.Context) error {
	return nil
}
