package db

type DeerDB interface {
	SaveLoginToken(userToken, loginToken string) error
	GetLoginToken(userToken string) (string, error)
	GetUserToken(loginToken string) (string, error)
	RemoveToken(userToken, loginToken string) error
}
