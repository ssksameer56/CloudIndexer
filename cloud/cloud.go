package cloud

import (
	"context"
	"time"
)

type Cloud interface {
	GetFiles(ctx context.Context, name string) ([]string, error)
	Connect(ctx context.Context) error
	CheckForChange(ctx context.Context, cursor string, timeout time.Duration, notifcationChannel <-chan bool)
	GetPointerToPath(ctx context.Context, path string) (string, error)
}
