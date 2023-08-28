package main

import (
	"fmt"
	"net/http"

	"gin-framework/src/auth"
	"gin-framework/src/db"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

var Usersinfor, err = db.GetDB()

func main() {
	r := setupRouter()

	// Listen and Serve on 0.0.0.0:8080
	r.Run(":8080")
}

func setupRouter() *gin.Engine {
	r := gin.Default()

	r.GET("/one", getting)
	r.GET("/two/:name", getName)
	r.POST("/login", postLogin)
	r.GET("/admin", Printoken)
	r.POST("/Register", func(c *gin.Context) {
		userRegister := auth.UserRegister{}
		userRegister.Register(c, Usersinfor)
	})
	PrivateInfo := r.Group("/Private")
	PrivateInfo.Use(auth.AuthLogin)
	{
		PrivateInfo.GET("/Friend")
		PrivateInfo.GET("/history")
	}
	r.GET("/Refresh", auth.Refreshtoken)
	return r
}

func getting(c *gin.Context) {
	c.String(http.StatusOK, "getting")
}

func getName(c *gin.Context) {
	name := c.Param("name")

	c.String(http.StatusOK, "GetName "+name)
}

func postLogin(c *gin.Context) {
	username := c.PostForm("username")
	password := c.PostForm("password")
	fmt.Println("Usersinfor: ", Usersinfor)
	for _, user := range Usersinfor {
		if user.Username == username {
			fmt.Println("Password:", user.Password)
			if bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)) == nil {
				token, err := auth.GenerateToken(string(user.ID))
				refreshToken, _ := auth.GenerateRefreshToken(user.ID)
				if err == nil {
					c.JSON(http.StatusOK, gin.H{"Token": token, "refreshtoken": refreshToken})
				}
				return
			}
		}
	}
	c.JSON(http.StatusInternalServerError, gin.H{"token": err})
	return
}

func Printoken(c *gin.Context) {
	c.JSON(http.StatusOK, auth.Store.GetAllTokens())
}
