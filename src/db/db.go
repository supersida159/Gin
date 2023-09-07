package db

import (
	"database/sql"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
)

type userInfor struct {
	ID          uint64
	Username    string
	Password    string
	Phonenumber string
}
type UsersinforDB []userInfor

var Usersinfor, Err = GetDB()

func UpdateDB() {
	Usersinfor, Err = GetDB()
	return
}

func ConnectSQL() (*sql.DB, error) {
	dsn := "root:Tungpro123@@tcp(localhost:3306)/chatting"
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

func GetDB() ([]userInfor, error) {
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

func AddRow(username, password, phoneNumber string) error {
	db, err := ConnectSQL()
	if err != nil {
		return err
	}
	defer db.Close()

	query := "INSERT INTO userinfor (UserName, Password, PhoneNumber) VALUES (?, ?, ?)"
	_, err = db.Exec(query, username, password, phoneNumber)
	return err
}
