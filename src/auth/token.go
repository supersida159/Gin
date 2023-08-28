package auth

import (
	"time"

	"fmt"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

var jwtKey = []byte("Golang-language")

type CustomClaims struct {
	UserID string `json:"user_id"`
	jwt.StandardClaims
}

var Store = NewTokenStore()

func GenerateToken(userID string) (string, error) {

	expirationTime := time.Now().Add(30 * time.Second) // Access token expiration time
	claims := &CustomClaims{
		UserID: userID,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	Token, err := token.SignedString(jwtKey)
	if err == nil {
		Store.StoreToken(Token, expirationTime, string(userID))
	}
	return token.SignedString(jwtKey)
}
func GenerateRefreshToken(userID uint64) (string, error) {
	expirationTime := time.Now().Add(5 * time.Minute)
	claims := &CustomClaims{
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	refreshtoken, err := token.SignedString(jwtKey)
	if err == nil {
		Store.StoreToken(refreshtoken, expirationTime, string(userID)+"R")
	}
	return token.SignedString(jwtKey)
}

func DecodeToken(TK string) (string, error) {
	claims := &CustomClaims{}
	_, err := jwt.ParseWithClaims(TK, claims, func(token *jwt.Token) (interface{}, error) {
		// ...
		return jwtKey, nil
	})
	return claims.UserID, err
}

func AuthLogin(c *gin.Context) {
	accessToken := c.GetHeader("Authorization")
	// refreshToken := c.GetHeader("X-Refresh-Token")
	for _, tokenString := range Store.tokens {
		if accessToken == tokenString.token {
			token, err := jwt.Parse(accessToken, func(token *jwt.Token) (interface{}, error) {
				// Provide the secret key used for token verification
				return jwtKey, nil
			})
			if err != nil {
				fmt.Println("Error parsing token:", err)
				return
			}
			if _, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
				// Access claims from the token

				c.JSON(200, gin.H{"token status": "Validtoken"})
				return
				// ... access other claims
			} else {
				c.JSON(200, gin.H{"token status": "expire accestoken"})
				return
			}
		}
	}
}

func Refreshtoken(c *gin.Context) {
	accessToken := c.GetHeader("Authorization")
	refreshToken := c.GetHeader("X-Refresh-Token")
	for userID, tokenString := range Store.tokens {
		if accessToken == tokenString.token {
			_, err := jwt.Parse(accessToken, func(token *jwt.Token) (interface{}, error) {
				// Provide the secret key used for token verification
				return jwtKey, nil
			})
			ve, ok := err.(*jwt.ValidationError)
			if ok && ve.Errors&jwt.ValidationErrorExpired != 0 {
				// check if access invalid
				if refreshToken == Store.tokens[userID+"R"].token {
					token, err := jwt.Parse(Store.tokens[userID+"R"].token, func(token *jwt.Token) (interface{}, error) {
						// Provide the secret key used for token verification
						return jwtKey, nil
					})
					if err != nil {
						fmt.Println("Error parsing token:", err)
						return
					}
					if _, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
						// Access claims from the token
						GenerateToken(userID)
						c.JSON(200, gin.H{"token status": "Update accestoken"})
						return
					}
					c.JSON(200, gin.H{"refreshtoken status": "invalid"})
					return
					// ... access other claims
				}
				c.JSON(200, gin.H{"refreshtoken status": "invalid"})
				return
			}
		}
	}
	c.JSON(200, gin.H{"token status": "not exist"})
}
