package auth

import (
	"time"

	"fmt"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

var jwtKey = []byte("Golang-language")

type CustomClaims struct {
	UserID int `json:"user_id"`
	jwt.StandardClaims
}

var Store = NewTokenStore()

func GenerateToken(userID int) (string, string, error) {

	accessExpirationTime := time.Now().Add(10 * time.Minute) // Access token expiration time
	refreshExpirationTime := time.Now().Add(10 * time.Hour)  // Refresh token expiration time

	// Generate the access token
	accessClaims := &CustomClaims{
		UserID: userID,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: accessExpirationTime.Unix(),
		},
	}
	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, accessClaims)
	access, err := accessToken.SignedString(jwtKey)
	if err != nil {
		return "", "", err
	}

	// Generate the refresh token
	refreshClaims := &CustomClaims{
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: refreshExpirationTime.Unix(),
		},
	}
	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims)
	refresh, err := refreshToken.SignedString(jwtKey)
	if err != nil {
		return "", "", err
	}

	// Store both tokens
	Store.StoreToken(access, refresh, accessExpirationTime, userID)
	fmt.Println(Store.Tokens[userID].Token)

	return access, refresh, nil
}

func DecodeToken(TK string) (int, error) {
	claims := &CustomClaims{}
	_, err := jwt.ParseWithClaims(TK, claims, func(token *jwt.Token) (interface{}, error) {
		// ...
		return jwtKey, nil
	})
	return claims.UserID, err
}

func AuthLogin(c *gin.Context) {

	accessToken := c.GetHeader("Authorization")
	if accessToken == "" {
		c.JSON(400, gin.H{"token status": "Missing Token"})
		c.AbortWithStatus(400)
		return

	}
	// refreshToken := c.GetHeader("X-Refresh-Token")
	for _, tokenString := range Store.Tokens {
		if accessToken == tokenString.Token {
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
				c.Next()
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
	for userID, tokenString := range Store.Tokens {
		if accessToken == tokenString.Token {
			_, err := jwt.Parse(accessToken, func(token *jwt.Token) (interface{}, error) {
				// Provide the secret key used for token verification
				return jwtKey, nil
			})
			ve, ok := err.(*jwt.ValidationError)
			if ok && ve.Errors&jwt.ValidationErrorExpired != 0 {
				// check if access invalid
				if refreshToken == Store.Tokens[userID].RefreshToken {
					token, err := jwt.Parse(Store.Tokens[userID].RefreshToken, func(token *jwt.Token) (interface{}, error) {
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
	c.JSON(401, gin.H{"token status": "not exist"})
	c.AbortWithStatus(401)
}
