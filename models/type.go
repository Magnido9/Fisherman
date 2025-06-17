package models

type UrlCheckResult struct {
	IsPhishing bool   `json:"isPhishing"`
	Error      string `json:"error,omitempty"`
}

type AsyncUrlCheckResult struct {
	Index  int
	Result UrlCheckResult
}
