package models_db

import (
	"log"

	"github.com/dgrijalva/jwt-go"
	_ "github.com/go-sql-driver/mysql"
)

type Claims struct {
	Login_user string `json:"login"`
	Role       string `json:"role"`
	jwt.StandardClaims
}

type User struct {
	Login_user    string `json:"login"`
	Password_user string `json:"password"`
	Role_user     string `json:"role"`
	First_name    string `json:"first_name"`
	Last_name     string `json:"last_name"`
	Date_both     string `json:"date_both"`
	Mobile_number string `json:"mobile_number"`
	Address       string `json:"address"`
}

func CheckUser(user *User) (bool, error, string) {
	rows, err := db.Query(`SELECT *
							from (
							select login, password, role 
							from clients 
							union all 
							select login, password, role 
							from employees) as Users
							where login = ? and password = ?`, &user.Login_user, &user.Password_user)
	if err != nil {
		log.Fatal(err)
		return false, err, "n"
	}
	defer rows.Close()

	var log string
	var pass string
	var role string

	for rows.Next() {
		err := rows.Scan(&log, &pass, &role)
		if err != nil {
			return false, err, "n"
		}
		if log == user.Login_user && pass == user.Password_user {
			return true, nil, role
		}
	}
	return false, err, ""

}
func AddUser(user *User) (err error) {
	_, err = db.Exec(`insert into Clients(login,password,first_name,last_name,date_both,mobile_number,address) values(?, ?, ?, ?, ?, ?, ?)`,
		&user.Login_user, &user.Password_user, &user.First_name, &user.Last_name, &user.Date_both, &user.Mobile_number, &user.Address)
	if err != nil {
		log.Fatal(err)
		return err
	}
	return nil
}
