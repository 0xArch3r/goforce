package types

type SearchResults struct {
	SearchRecords []SObject `json:"searchRecords"`
}

type QueryResult struct {
	TotalSize      int       `json:"totalSize"`
	Done           bool      `json:"done"`
	NextRecordsURL string    `json:"nextRecordsUrl"`
	Records        []SObject `json:"records"`
}
