package workers

import (
	"github.com/ssksameer56/CloudIndexer/cloud"
)

type CloudWatcher struct {
	CloudProvider cloud.Cloud
	AuthCode      string
}

func (cw *CloudWatcher) WatchCloud() {

}
