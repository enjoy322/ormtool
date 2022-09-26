package ormtool

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"log"
)

//mysql connect
func mysqlConn(conn string) *sql.DB {
	db, err := sql.Open("mysql", conn)
	if err != nil {
		log.Fatalln(err)
	}
	err = db.Ping()
	if err != nil {
		log.Fatalln(err)
	}
	log.Println("MySQL connected")
	return db
}
