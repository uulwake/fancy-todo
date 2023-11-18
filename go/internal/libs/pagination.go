package libs

type Pagination struct {
	Size int `json:"size,omitempty"`
	Number int `json:"number,omitempty"`
	Total int64 `json:"total,omitempty"`
}