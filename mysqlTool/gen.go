package mysqlTool

import (
	"database/sql"
	"ormtool/base"
)

func GenMySQL(my base.MysqlConfig, c base.Config) {
	db := dbConn(my)
	defer func(db *sql.DB) {
		err := db.Close()
		if err != nil {
			panic(err)
		}
	}(db)
	packageName, fileDir, fileName, data := Service(db).StructContent(my.Database, c)
	base.Write(packageName, fileDir, fileName, data, c.IsGenInOneFile)
}
