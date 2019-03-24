package constants

const (
	AuthAction            = "authenticate"
	CreateAction          = "create"
	DeleteAction          = "delete"
	UpdateAction          = "update"
	CollectAction         = "collect"
	CollectSecondHWAction = "collectHW"

	StatusError         = "ERROR"
	StatusSuccess       = "SUCCESS"
	StatusAuthenticated = "AUTHENTICATED"
	StatusCreated       = "CREATED"
	StatusConflict      = "CONFLICT"

	NOT_MATCH             = "password not match"
	ACCOUNT_ALREADY_EXIST = "account already exist"
	ACCOUNT_NOT_FOUND     = "account not found"
)
