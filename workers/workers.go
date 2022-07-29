package workers

import (
	"context"
	"sync"
)

type Worker interface {
	Init(ctx context.Context) error
	Run(wg *sync.WaitGroup)
	Stop() error
}

//TODO : workers loop write to call from main.go
