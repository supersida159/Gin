package Message

import (
	"fmt"
	"gin-framework/src/auth"
	"gin-framework/src/db"
	"net/http"
	"time"

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

func (m *Mess) GetMessFirst(c *gin.Context) {
	type RequestData struct {
		NumberScroll int `json:"Number"`
	}
	var requestData RequestData
	type ChattingHisWUser struct {
		ChattingID int
		Sender     int
		Receiver   int
		Content    string
		Sendtime   time.Time
		UserName   string
	}

	var Result []ChattingHisWUser
	var UniqueHis []int

	fmt.Println(c.GetHeader("Authorization"))

	Sender, _ := auth.DecodeToken(c.GetHeader("Authorization"))
	if err := c.ShouldBindJSON(&requestData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"Status": err})
		return
	}
	fmt.Println(requestData.NumberScroll)
	for _, Messx := range db.ChattingHisDB {
		if Sender == Messx.Sender {
			uExists := false
			for _, Histlist := range UniqueHis {
				if Messx.Receiver == Histlist {
					uExists = true
					break
				}

			}
			if !uExists {
				UniqueHis = append(UniqueHis, Messx.Receiver)
				if len(Result) < requestData.NumberScroll*10 {
					chattingHisWUser := ChattingHisWUser{
						ChattingID: Messx.ChattingID,
						Sender:     Messx.Sender,
						Receiver:   Messx.Receiver,
						Content:    Messx.Content,
						Sendtime:   Messx.Sendtime,            // Copy the ChattingHis data
						UserName:   db.IDtoName[Messx.Sender], // Replace with the actual username
					}
					Result = append(Result, chattingHisWUser)
				}
			}
		} else if Sender == Messx.Receiver {
			uExists := false
			for _, Histlist := range UniqueHis {
				if Messx.Sender == Histlist {
					uExists = true
					break
				}
			}
			if !uExists {
				UniqueHis = append(UniqueHis, Messx.Sender)

				if len(Result) < requestData.NumberScroll*10 {
					chattingHisWUser := ChattingHisWUser{
						ChattingID: Messx.ChattingID,
						Sender:     Messx.Sender,
						Receiver:   Messx.Receiver,
						Content:    Messx.Content,
						Sendtime:   Messx.Sendtime,              // Copy the ChattingHis data
						UserName:   db.IDtoName[Messx.Receiver], // This is last user send message
					}
					Result = append(Result, chattingHisWUser)
				}
			}
		}

	}
	c.JSON(http.StatusOK, Result)

}
