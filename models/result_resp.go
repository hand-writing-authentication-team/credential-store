package models

type ResultResp struct {
	JobID     string `json:"jobid"`
	Status    string `json:"status"`
	ErrorMsg  string `json:"error_message,omitempty"`
	TimeStamp int64  `json:"timestamp"`
}
