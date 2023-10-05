package main

import (
	"net/http"

	"gin-framework/src/Message"
	"gin-framework/src/auth"
	"gin-framework/src/db"
	"gin-framework/src/friend"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

func main() {
	r := setupRouter()

	// Listen and Serve on 0.0.0.0:8080
	r.Run(":8080")
}

func setupRouter() *gin.Engine {
	r := gin.Default()

	// Enable CORS for all origins
	r.Use(corsMiddleware())
	r.GET("/one", getting)
	r.GET("/two/:name", getName)
	r.POST("/login", postLogin)
	r.GET("/admin", Printoken)
	r.POST("/Register", func(c *gin.Context) {
		userRegister := auth.UserRegister{}
		userRegister.Register(c)
	})
	PrivateInfo := r.Group("/Private")
	PrivateInfo.Use(auth.AuthLogin)
	{
		PrivateInfo.GET("/Friend", friend.NewFriendService().GetFriendshipDB)
		PrivateInfo.GET("/History", Message.NewMess().GetMess)
		PrivateInfo.GET("/SendMess", Message.NewMess().SendMess)
	}
	//friend...
	//seaching friend
	r.POST("/SeachFriend", auth.AuthLogin, friend.NewFriendService().Search)
	r.POST("/UpdateFriend", auth.AuthLogin, friend.NewFriendService().ReceiveUpdateRelate)
	// verify refreshtoken
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
	if db.UsersinforDB == nil {
		c.JSON(http.StatusInternalServerError, gin.H{"Error :": "DB error"})
	}
	for _, user := range db.UsersinforDB {
		if user.Username == username {
			if bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)) == nil {
				_, _, err := auth.GenerateToken(user.ID)
				if err == nil {
					c.JSON(http.StatusOK, auth.Store.Tokens[user.ID])
					return
				} else {
					c.JSON(http.StatusInternalServerError, gin.H{"GenerateRefreshToken err": err})
					return
				}

			}

		}
	}
	c.JSON(http.StatusInternalServerError, gin.H{"token": db.Err})
}

func Printoken(c *gin.Context) {
	c.JSON(http.StatusOK, auth.Store.GetAllTokens())
}
func corsMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*") // You may want to restrict this to specific domains in production
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Authorization, Content-Type")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}

		c.Next()
	}
}
