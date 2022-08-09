package elasticservice

import (
	"encoding/base64"

	"github.com/ssksameer56/CloudIndexer/models"
)

func Hash(data models.TextStoreModel) string {
	key := base64.StdEncoding.EncodeToString([]byte(data.FilePath + data.Name))
	return key
}
