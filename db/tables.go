package db

const TokenTable = "t_token_info"

type UserTokenInfo struct {
	UserToken  string
	LoginToken string
}
