package db

import (
	"database/sql"
	"fmt"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

type userInfor struct {
	ID          int
	Username    string
	Password    string
	Phonenumber string
}

type Friendship struct {
	FriendshipID int
	User1ID      int
	User2ID      int
	Status       string
	Sender       int
}
type ChattingHis struct {
	ChattingID int
	Sender     int
	Receiver   int
	Content    string
	Sendtime   time.Time
}

var FriendshipDB, ErrF = GetFriendshipDB()

func UpdateFriendshipDB() {
	FriendshipDB, ErrF = GetFriendshipDB()
}

var UsersinforDB, Err = GetUserDB()

func UpdateUserDB() {
	UsersinforDB, Err = GetUserDB()
}

var ChattingHisDB, ErrC = GetMessage()

func UpdateChattingHis() {
	ChattingHisDB, ErrC = GetMessage()
}

func ConnectSQL() (*sql.DB, error) {
	dsn := "root:Tungpro123@@tcp(localhost:3306)/chatting?parseTime=true"
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}
	fmt.Print("connect success")
	return db, nil
}

func GetUserDB() ([]userInfor, error) {
	db, err := ConnectSQL()
	if err != nil {
		fmt.Println("connect error")
		return nil, err
	}
	defer db.Close()

	rows, err := db.Query("SELECT * FROM userinfor")
	if err != nil {
		fmt.Println("select table error")
		return nil, err
	}
	defer rows.Close()

	var UsersinforDB []userInfor

	for rows.Next() {
		var userDB userInfor
		err := rows.Scan(&userDB.ID, &userDB.Username, &userDB.Phonenumber, &userDB.Password)
		if err != nil {
			fmt.Println("get row error")
			return nil, err
		}

		UsersinforDB = append(UsersinforDB, userDB)
	}
	if err := rows.Err(); err != nil {

		return nil, err

	}
	return UsersinforDB, nil
}

// get Friendship
func GetFriendshipDB() ([]Friendship, error) {
	db, err := ConnectSQL()
	if err != nil {
		fmt.Println("connect error")
		return nil, err
	}
	defer db.Close()

	rows, err := db.Query("SELECT * FROM friendship")
	if err != nil {
		fmt.Println("select table error")
		return nil, err
	}
	defer rows.Close()

	var FriendShipDB []Friendship

	for rows.Next() {
		var Frienddb Friendship
		err := rows.Scan(&Frienddb.FriendshipID, &Frienddb.User1ID, &Frienddb.User2ID, &Frienddb.Status, &Frienddb.Sender)
		if err != nil {
			fmt.Println("get row error")
			return nil, err
		}

		FriendShipDB = append(FriendShipDB, Frienddb)
	}
	if err := rows.Err(); err != nil {

		return nil, err

	}
	return FriendShipDB, nil
}

func AddRowUser(username, password, phoneNumber string) error {
	db, err := ConnectSQL()
	if err != nil {
		return err
	}
	defer db.Close()

	query := "INSERT INTO userinfor (UserName, Password, PhoneNumber) VALUES (?, ?, ?)"
	_, err = db.Exec(query, username, password, phoneNumber)
	return err
}

func AddRowFriendship(User1ID, User2ID int, Status string) error {
	db, err := ConnectSQL()
	if err != nil {
		return err
	}
	defer db.Close()

	query := "INSERT INTO friendship (User1ID, User2ID, Status,Sender) VALUES (?, ?, ?, ?)"
	_, err = db.Exec(query, User1ID, User2ID, Status, User1ID)
	UpdateFriendshipDB()
	return err
}
func SendmessageDB(Sender, Receiver int, Content string, ChattingID int) error {
	db, err := ConnectSQL()
	if err != nil {
		return err
	}
	defer db.Close()

	query := "INSERT INTO chatting_history (Sender, Receiver, Content,ChattingID) VALUES (?, ?, ?, ?)"
	_, err = db.Exec(query, Sender, Receiver, Content, ChattingID)
	UpdateChattingHis()
	return err

}
func GetMessage() ([]ChattingHis, error) {
	db, err := ConnectSQL()
	if err != nil {
		fmt.Println("connect error")
		return nil, err
	}
	defer db.Close()

	rows, err := db.Query("SELECT * FROM chatting_history")
	if err != nil {
		fmt.Println("select table error")
		return nil, err
	}
	defer rows.Close()

	var ChattingHisDB []ChattingHis

	for rows.Next() {
		var ChattingDB ChattingHis
		err := rows.Scan(&ChattingDB.ChattingID, &ChattingDB.Sender, &ChattingDB.Receiver, &ChattingDB.Content, &ChattingDB.Sendtime)
		if err != nil {
			fmt.Println("get row error")
			fmt.Println(err)
			return nil, err
		}

		ChattingHisDB = append(ChattingHisDB, ChattingDB)
	}
	if err := rows.Err(); err != nil {

		return nil, err

	}
	return ChattingHisDB, nil

}
func UpdateFriendship(Status string, Sender, FriendshipID int) error {
	db, err := ConnectSQL()
	if err != nil {
		return err
	}
	defer db.Close()

	query := "UPDATE friendship SET Status=? ,Sender=? WHERE FriendshipID=?"
	_, err = db.Exec(query, Status, Sender, FriendshipID)
	UpdateFriendshipDB()
	return err
}
