package friend

import (
	"gin-framework/src/db"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

var SearchingFriendResult map[string]string

func SeachingFriend(c *gin.Context) {
	searching := c.PostForm("searchingUD")
	_, err := strconv.Atoi(searching) //check input is num or username
	if err == nil {
		for _, user := range db.Usersinfor {
			if strings.Contains(user.Phonenumber, searching) {
				SearchingFriendResult[user.Username] = user.Phonenumber
			}
		}
	} else {
		for _, user := range db.Usersinfor {
			if strings.Contains(searching, user.Username) {
				SearchingFriendResult[user.Username] = user.Phonenumber
			}
		}
	}
	if len(SearchingFriendResult) > 10 {
		c.JSON(http.StatusBadRequest, gin.H{"ERR status": "Result leng over 10"})
	} else if len(SearchingFriendResult) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"ERR status": "no user"})
	} else {
		println("somethingwrong")
		c.JSON(http.StatusOK, SearchingFriendResult)
	}
}
