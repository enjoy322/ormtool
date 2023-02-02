package ormtool

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
)

// mysql connect
func mysqlConn(conn string) (*sql.DB, error) {
	db, err := sql.Open("mysql", conn)
	if err != nil {
		return nil, err
	}
	err = db.Ping()
	if err != nil {
		return nil, err
	}
	return db, nil
}
