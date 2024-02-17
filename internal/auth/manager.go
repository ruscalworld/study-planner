package auth

type Manager interface {
	Authenticate(userInfo *UserInfo) (*Token, error)
	Authorize(token *Token) (*TokenInfo, error)
	GetTokenProvider() TokenProvider
}
