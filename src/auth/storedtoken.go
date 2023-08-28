package auth

import (
	"sync"
	"time"
)

type tokenStruct struct {
	expirationdate time.Time
	token          string
}
type TokenStore struct {
	mu     sync.Mutex
	tokens map[string]tokenStruct
}

func NewTokenStore() *TokenStore {
	return &TokenStore{
		tokens: make(map[string]tokenStruct),
	}
}

func (ts *TokenStore) StoreToken(token string, expiration time.Time, UserID string) {
	ts.mu.Lock()
	defer ts.mu.Unlock()
	ts.tokens[string(UserID)] = tokenStruct{
		expirationdate: expiration,
		token:          token,
	}
}

func (ts *TokenStore) IsTokenValid(token string) bool {
	ts.mu.Lock()
	defer ts.mu.Unlock()
	for _, tokenData := range ts.tokens {
		if tokenData.token == token {
			return tokenData.expirationdate.After(time.Now())
		}
	}
	return false
}
func (ts *TokenStore) GetAllTokens() map[string]time.Time {
	ts.mu.Lock()
	defer ts.mu.Unlock()
	tokensCopy := make(map[string]time.Time)
	for _, TokenData := range ts.tokens {
		tokensCopy[TokenData.token] = TokenData.expirationdate
	}
	return tokensCopy
}
