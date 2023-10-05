package auth

import (
	"sync"
	"time"
)

type tokenStruct struct {
	expirationdate time.Time
	Token          string
	RefreshToken   string
}
type TokenStore struct {
	mu     sync.Mutex
	Tokens map[int]tokenStruct
}

func NewTokenStore() *TokenStore {
	return &TokenStore{
		Tokens: make(map[int]tokenStruct),
	}
}

func (ts *TokenStore) StoreToken(token, refreshtoken string, expiration time.Time, UserID int) {
	ts.mu.Lock()
	defer ts.mu.Unlock()
	ts.Tokens[UserID] = tokenStruct{
		expirationdate: expiration,
		Token:          token,
		RefreshToken:   refreshtoken,
	}
}

func (ts *TokenStore) IsTokenValid(token string) bool {
	ts.mu.Lock()
	defer ts.mu.Unlock()
	for _, tokenData := range ts.Tokens {
		if tokenData.Token == token {
			return tokenData.expirationdate.After(time.Now())
		}
	}
	return false
}
func (ts *TokenStore) GetAllTokens() map[string]time.Time {
	ts.mu.Lock()
	defer ts.mu.Unlock()
	tokensCopy := make(map[string]time.Time)
	for _, TokenData := range ts.Tokens {
		tokensCopy[TokenData.Token] = TokenData.expirationdate
	}
	return tokensCopy
}
