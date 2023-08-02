package main

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
)

var Usersinfor, err = connectSQL()

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
			if user.Password == password {
				c.JSON(http.StatusOK, gin.H{"token": "this is a token"})
				return
			}
		}
	}
	c.JSON(http.StatusInternalServerError, gin.H{"token": "wrong pass word or username"})
	return
}
