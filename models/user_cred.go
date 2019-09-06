package models

type UserCredentials struct {
	ID              string
	Username        string
	PasswordContent string
	Handwriting     string
	UserModel       string
	Race            string
	Created         int64
	Modified        int64
	Deleted         bool
}

type UserValidateHW struct {
	ID          string
	UserID      string
	Username    string
	Handwriting string
	Created     int64
	Modified    int64
	Deleted     bool
}
