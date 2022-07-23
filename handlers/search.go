package handlers

import (
	"context"

	"github.com/ssksameer56/CloudIndexer/cloud"
)

type SearchHandler struct {
	Service         cloud.Cloud
	ESSearchService interface{}
}

func (sh *SearchHandler) SearchText(ctx context.Context, keyword string) ([]string, error) {
	return []string{}, nil
}
