package models_db

import (
	"database/sql"
	"log"

	_ "github.com/go-sql-driver/mysql"
)

var db *sql.DB

func InitDB(dataSourceName string) {
	var err error
	db, err = sql.Open("mysql", dataSourceName)
	if err != nil {
		log.Panic(err)
		return
	}
	err = db.Ping()
	if err != nil {
		log.Panic(err)
		return
	}
}
