package models

type Feature struct {
	DataPoints map[string]interface{} `json:"data_points"`
}

type FeatureReq struct {
	UserHandwriting string    `json:"user_handwriting"`
	PrevDataPoints  []Feature `json:"previous_data_points,omitempty"`
}

type ValidateStatus struct {
	Status bool `json:"status"`
}
