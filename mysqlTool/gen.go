package mysqlTool

import (
	"database/sql"
	"log"

	"github.com/enjoy322/ormtool/base"
)

func GenMySQL(my base.MysqlConfig, c base.Config) {
	db := dbConn(my)
	defer func(db *sql.DB) {
		err := db.Close()
		if err != nil {
			log.Fatalln(err)
		}
	}(db)
	fileData, data := Service(db, my.Database, c).GenStruct()
	// write into file
	base.Write(fileData, data, c.IsGenInOneFile)
}
