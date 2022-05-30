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
	fielData, data := Service(db).GenStruct(my.Database, c)
	// write into file
	base.Write(fielData, data, c.IsGenInOneFile)
}
