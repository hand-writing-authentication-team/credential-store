package models

type AuthenticationRequest struct {
	JobID     string `json:"jobid"`
	Username  string `json:"username"`
	Password  string `json:"password"`
	Handwring string `json:"handwriting"`
	Race      string `json:"race,omitempty"`
	Action    string `json:"action"`
}
