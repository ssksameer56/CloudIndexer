package models

type ESErrorResponse struct {
	Info *ESErrorInfo `json:"error,omitempty"`
}

type ESErrorInfo struct {
	RootCause []*ESErrorInfo
	Type      string
	Reason    string
	Phase     string
}

type ESIndexResponse struct {
	Index   string `json:"_index"`
	ID      string `json:"_id"`
	Version int    `json:"_version"`
	Result  string
}

type ESSearchResponse struct {
	Took int64
	Hits struct {
		Total struct {
			Value int64
		}
		Hits []*SearchHit
	}
}

type SearchHit struct {
	Score   float64 `json:"_score"`
	Index   string  `json:"_index"`
	Type    string  `json:"_type"`
	Version int64   `json:"_version,omitempty"`

	Source TextStoreModel `json:"_source"` // the struct containing your data
}

type Match struct {
	Text string `json:"text,omitempty"`
}
type Query struct {
	Match Match `json:"match,omitempty"`
}
type SearchRequest struct {
	Query `json:"query,omitempty"`
}
