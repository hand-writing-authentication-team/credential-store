package models

type AuthenticationRequest struct {
	JobID     string `json:"jobid"`
	Username  string `json:"username"`
	Password  string `json:"password"`
	Handwring string `json:"handwriting"`
	Action    string `json:"action"`
}
