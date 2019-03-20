package models

type UserCredentials struct {
	ID              string
	Username        string
	PasswordContent string
	Handwriting     string
	Race            string
	Created         int64
	Modified        int64
	Deleted         bool
}
