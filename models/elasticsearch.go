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
	Index   string `json:"_index,omitempty"`
	ID      string `json:"_id,omitempty"`
	Version int    `json:"_version,omitempty"`
	Result  string `json:"result,omitempty"`
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

type ESGetResponse struct {
	Index       string `json:"_index"`
	ID          string `json:"_id"`
	Version     int    `json:"_version"`
	SeqNo       int    `json:"_seq_no"`
	PrimaryTerm int    `json:"_primary_term"`
	Found       bool   `json:"found"`
}

type SearchHit struct {
	Score   float64 `json:"_score,omitempty"`
	Index   string  `json:"_index,omitempty"`
	Type    string  `json:"_type,omitempty"`
	Id      string  `json:"_id,omitempty"`
	Version int64   `json:"_version,omitempty"`

	Source TextStoreModel `json:"_source,omitempty"` // the struct containing your data
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
