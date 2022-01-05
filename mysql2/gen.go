package mysql2

import (
	"convert/base"
	"fmt"
	"log"
)

func GenMySQL(c base.Config) {
	log.Println("-----generate-----")
	db := dbConn(c.MySQL)
	if c.IsGenInOneFile {
		Service(db).CreateOneFile(c)
	} else {
		fmt.Println("fenli")
		Service(db).Create(c)
	}
}
