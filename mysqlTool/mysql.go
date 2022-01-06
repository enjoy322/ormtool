package mysqlTool

import (
	"convert/base"
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"strings"
)

//mysql连接
func dbConn(c base.MysqlConfig) *sql.DB {
	conn := strings.Builder{}
	conn.WriteString(c.User)
	conn.WriteString(":")
	conn.WriteString(c.Password)
	conn.WriteString("@tcp(")
	conn.WriteString(c.Host)
	conn.WriteString(":")
	conn.WriteString(c.Port)
	conn.WriteString(")/")
	conn.WriteString(c.Database)
	conn.WriteString("?parseTime=true")
	db, err := sql.Open("mysql", conn.String())
	if err != nil {
		panic(err)
	}
	err = db.Ping()
	if err != nil {
		panic(err)
	}
	fmt.Println("MySQL connected")
	return db
}
