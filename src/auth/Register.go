package auth

import (
	"net/http"
	"strconv"

	"gin-framework/src/db"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	// "gin-framework/src/db"
)

type UsersinforDB db.UsersinforDB
type UserRegister struct{}

func (au UserRegister) Register(c *gin.Context, data UsersinforDB) {
	username := c.PostForm("username")
	password := c.PostForm("password")
	Phonenumber := c.PostForm("phonenumber")
	if Phonenumber == "" {
		//Add the phonenumber to empty(use the last phonenumber to make it)
		NumPhone, _ := strconv.Atoi(data[len(data)-1].Phonenumber)
		NumPhone = NumPhone + 101
		Phonenumber = "0" + strconv.Itoa(NumPhone)
	}
	for _, user := range data {
		if user.Username == username {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Username already taken"})
			return
		}
		if user.Phonenumber == username {
			c.JSON(http.StatusBadRequest, gin.H{"error": "phonenumer already taken"})
			return
		}
	}
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not hash password"})
		return
	}
	if db.AddRow(username, string(hashedPassword), Phonenumber) == nil {
		c.JSON(http.StatusOK, gin.H{"message": "User registered successfully"})
		db.UpdateDB()
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "User registered fail", "error": db.AddRow(username, string(hashedPassword), Phonenumber)})
}
