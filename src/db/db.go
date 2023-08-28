package db

import (
	"database/sql"

	_ "github.com/go-sql-driver/mysql"
)

type userInfor struct {
	ID          uint64
	Username    string
	Password    string
	Phonenumber string
	Token       any
}
type UsersinforDB []userInfor

func ConnectSQL() (*sql.DB, error) {
	dsn := "root:Tungpro123@@tcp(localhost:3306)/chatting"
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}
	return db, nil
}

func GetDB() ([]userInfor, error) {
	db, err := ConnectSQL()
	if err != nil {
		return nil, err
	}
	defer db.Close()

	rows, err := db.Query("SELECT * FROM userinfor")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var UsersinforDB []userInfor

	for rows.Next() {
		var userDB userInfor
		err := rows.Scan(&userDB.Username, &userDB.Password, &userDB.ID, &userDB.Phonenumber, &userDB.Token)
		if err != nil {
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

	query := "INSERT INTO userinfor (name, password, phone_number) VALUES (?, ?, ?)"
	_, err = db.Exec(query, username, password, phoneNumber)
	return err
}
