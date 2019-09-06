package models

type Feature struct {
	UserModel string `json:"user_model"`
}

type FeatureReq struct {
	UserHandwriting string `json:"handwriting"`
	UserModel       string `json:"user_model,omitempty"`
}

type ValidateStatus struct {
	Status bool `json:"status"`
}
