package Message

import (
	"fmt"
	"gin-framework/src/auth"
	"gin-framework/src/db"
	"net/http"

	"github.com/gin-gonic/gin"
)

type Mess struct{}

func NewMess() *Mess {
	return &Mess{}
}

func (m *Mess) SendMess(c *gin.Context) {
	var MessSend struct {
		Message  string `json:"Message"`
		Receiver int    `json:"Receiver"`
	}
	Sender, _ := auth.DecodeToken(c.GetHeader("Authorization"))
	if err := c.ShouldBindJSON(&MessSend); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"Status": err})
		return
	}
	err := db.SendmessageDB(Sender, MessSend.Receiver, MessSend.Message, 0)
	if err == nil {
		c.JSON(http.StatusOK, gin.H{"Message sent": MessSend.Message})
	} else {
		fmt.Println("Error sending message:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
	}
}

func (m *Mess) GetMess(c *gin.Context) {
	var MessHis struct {
		Receiver int `json:"Receiver"`
	}
	var Result []db.ChattingHis
	Sender, _ := auth.DecodeToken(c.GetHeader("Authorization"))
	if err := c.ShouldBindJSON(&MessHis); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"Status": err})
		return
	}
	for _, Messx := range db.ChattingHisDB {
		if (Sender == Messx.Sender && MessHis.Receiver == Messx.Receiver) || (MessHis.Receiver == Messx.Sender && Sender == Messx.Receiver) {
			Result = append(Result, Messx)
		}

	}
	c.JSON(http.StatusOK, Result)
}
