package models_db

import (
	"log"

	_ "github.com/go-sql-driver/mysql"
)

type User struct {
	Login_user    string `json:"login"`
	Password_user string `json:"password"`
}

func CheckUser(user *User) (bool, error) {
	rows, err := db.Query(`select login, password from Clients where login = ?`, user.Login_user)
	if err != nil {
		log.Fatal(err)
		return false, err
	}
	defer rows.Close()

	var log string
	var pass string

	for rows.Next() {
		err := rows.Scan(&log, &pass)
		if err != nil {
			return false, err
		}
		if log == user.Login_user && pass == user.Password_user {
			return true, nil
		}
	}
	return false, err

}
func AddUser(user *User) (err error) {
	_, err = db.Exec(`insert into Clients(login,password) values(?, ?)`, &user.Login_user, &user.Password_user)
	if err != nil {
		log.Fatal(err)
		return err
	}
	return nil
}
