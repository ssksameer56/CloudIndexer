package cloud

import "context"

type Cloud interface {
	GetFiles(ctx context.Context, name string) ([]string, error)
	Connect(ctx context.Context) error
}
