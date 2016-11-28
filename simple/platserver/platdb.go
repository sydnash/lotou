package main

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
)

var db *sql.DB

func OpenDB() {
	var err error
	db, err = sql.Open("mysql", "root:root@tcp(127.0.0.1:3306)/niuniu")
	if err != nil {
		fmt.Println("connect sql:", err)
		return
	}
}
