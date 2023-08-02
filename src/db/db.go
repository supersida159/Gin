package db

import (
	"database/sql"

	_ "github.com/go-sql-driver/mysql"
)

type userInfor struct {
	ID          int
	Username    string
	Password    string
	Phonenumber string
	Token       any
}

func connectSQL() ([]userInfor, error) {
	var UsersinforDB []userInfor
	dsn := "root:Tungpro123@@tcp(localhost:3306)/chatting"

	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return UsersinforDB, err
	}

	// Test the connection
	if err := db.Ping(); err != nil {
		return UsersinforDB, err
	}
	rows, err := db.Query("SELECT * FROM userinfor")
	defer db.Close()
	for rows.Next() {
		var userDB userInfor
		err := rows.Scan(&userDB.Username, &userDB.Password, &userDB.ID, &userDB.Phonenumber, &userDB.Token)
		if err != nil {
			return UsersinforDB, err
		}
		UsersinforDB = append(UsersinforDB, userDB)
	}
	if err := rows.Err(); err != nil {
		return UsersinforDB, err
	}
	return UsersinforDB, err
}
