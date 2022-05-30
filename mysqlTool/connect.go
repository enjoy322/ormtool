package mysqlTool

import (
	"database/sql"
	"github.com/enjoy322/ormtool/base"
	_ "github.com/go-sql-driver/mysql"
	"log"
	"strings"
)

//mysql connect
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
		log.Fatalln(err)
	}
	err = db.Ping()
	if err != nil {
		log.Fatalln(err)
	}
	log.Println("MySQL connected")
	return db
}
