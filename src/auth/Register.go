package auth

import (
	"net/http"
	"strconv"

	"gin-framework/src/db"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	// "gin-framework/src/db"
)

var x, err, _ = db.GetUserDB()

type UserRegister struct{}

func (au UserRegister) Register(c *gin.Context) {
	username := c.PostForm("username")
	password := c.PostForm("password")
	Phonenumber := c.PostForm("phonenumber")
	if Phonenumber == "" {
		//Add the phonenumber to empty(use the last phonenumber to make it)
		NumPhone, _ := strconv.Atoi(db.UsersinforDB[len(db.UsersinforDB)-1].Phonenumber)
		NumPhone = NumPhone + 101
		Phonenumber = "0" + strconv.Itoa(NumPhone)
	}
	for _, user := range db.UsersinforDB {
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
	if db.AddRowUser(username, string(hashedPassword), Phonenumber) == nil {
		c.JSON(http.StatusOK, gin.H{"message": "User registered successfully"})
		db.UpdateUserDB()
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "User registered fail", "error": db.AddRowUser(username, string(hashedPassword), Phonenumber)})
}
