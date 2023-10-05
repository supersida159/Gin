package friend

import (
	"fmt"
	"gin-framework/src/auth"
	"gin-framework/src/db"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

type FriendService struct{}

func NewFriendService() *FriendService {
	return &FriendService{}
}

func (fs *FriendService) Search(c *gin.Context) {
	var searchInput struct {
		SearchingUD string `json:"searchingUD"`
	}

	if err := c.ShouldBindJSON(&searchInput); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	searching := searchInput.SearchingUD
	searchResults := fs.performSearch(searching)

	if len(searchResults) > 10 {
		c.JSON(http.StatusBadRequest, gin.H{"status": "Over 10 user", "data": searchResults})
	} else if len(searchResults) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"status": "No user"})
	} else {
		c.JSON(http.StatusOK, gin.H{"status": "Success", "data": searchResults})
	}
}

func (fs *FriendService) performSearch(searching string) map[string]string {
	searchResults := make(map[string]string)

	// Check if searching is numeric
	if _, err := strconv.Atoi(searching); err == nil {
		for _, user := range db.UsersinforDB {
			if strings.Contains(user.Phonenumber, searching) {
				searchResults[user.Username] = user.Phonenumber
			}
		}
	} else {
		for _, user := range db.UsersinforDB {
			if strings.Contains(user.Username, searching) {
				searchResults[user.Username] = user.Phonenumber
			}
		}
	}

	return searchResults
}
func (fs *FriendService) GetFriendshipDB(c *gin.Context) {
	var FSResult []db.Friendship
	UserID, _ := auth.DecodeToken(c.GetHeader("Authorization"))

	if len(db.FriendshipDB) < 1 {
		c.JSON(http.StatusOK, gin.H{"Status": "Empty FRS"})
		return
	}
	for _, friendship := range db.FriendshipDB {
		if friendship.User1ID == UserID || friendship.User2ID == UserID {
			FSResult = append(FSResult, friendship)
		}
	}

	c.JSON(http.StatusOK, FSResult)

}
func (fs *FriendService) ReceiveUpdateRelate(c *gin.Context) {
	var ReceiveRelate struct {
		Status string `json:"Status"`
		UserID int    `json:"UserID"`
	}
	UserID, _ := auth.DecodeToken(c.GetHeader("Authorization"))
	if err := c.ShouldBindJSON(&ReceiveRelate); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"Status": err})
		return
	}
	for _, friendship := range db.FriendshipDB {
		if (UserID == friendship.User1ID && ReceiveRelate.UserID == friendship.User2ID) || (UserID == friendship.User2ID && ReceiveRelate.UserID == friendship.User1ID) {
			fmt.Println("This relationship is exist")
			switch ReceiveRelate.Status {
			case "BL":
				if friendship.Status == "BL" {
					c.JSON(http.StatusConflict, gin.H{"Status": "BL"})
				} else {

					db.UpdateFriendship(ReceiveRelate.Status, UserID, friendship.FriendshipID)
					c.JSON(http.StatusConflict, gin.H{"Status": "Update BL"})
				}
			case "AD":
				if friendship.Status == "AD" || friendship.Status == "AC" {
					if UserID == friendship.Sender || friendship.Status == "AC" {
						c.JSON(http.StatusConflict, gin.H{"Status": "already AD"})
					} else {
						db.UpdateFriendship("AC", UserID, friendship.FriendshipID)
						c.JSON(http.StatusConflict, gin.H{"Status": "AC"})
					}
				} else if friendship.Status == "BL" {
					if friendship.Sender == UserID {
						db.UpdateFriendship("AD", UserID, friendship.FriendshipID)
						c.JSON(http.StatusOK, gin.H{"Status": "UB and AD"})
					} else {
						c.JSON(http.StatusConflict, gin.H{"Status": "you are Blocked"})
					}
				} else if friendship.Status == "UB" {
					db.UpdateFriendship("AD", UserID, friendship.FriendshipID)
				}
			case "UN":
				if friendship.Status == "AC" {
					db.UpdateFriendship("AD", ReceiveRelate.UserID, friendship.FriendshipID)
					c.JSON(http.StatusConflict, gin.H{"Status": "Un Friend"})
				} else {
					c.JSON(http.StatusConflict, gin.H{"Status": "you're not friend"})
				}
			case "UB":
				if friendship.Status == "BL" && friendship.Sender == UserID {
					db.UpdateFriendship("UB", UserID, friendship.FriendshipID)

					c.JSON(http.StatusConflict, gin.H{"Status": "Un block"})
				} else if friendship.Status == "BL" && friendship.Sender == ReceiveRelate.UserID {
					c.JSON(http.StatusConflict, gin.H{"Status": "you are block by other"})
				} else {
					c.JSON(http.StatusConflict, gin.H{"Status": "you are not block"})
				}
			}
			return
		}
	}
	err := db.AddRowFriendship(UserID, ReceiveRelate.UserID, ReceiveRelate.Status)
	if err == nil {
		c.JSON(http.StatusOK, gin.H{"Status": "Add new"})
	}
}
