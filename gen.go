package convert

import (
	"convert/base"
	"convert/mysql2"
	"log"
)

func Gen(c base.Config) {

	switch c.DataBaseType {
	case base.CodeDBMySQL:
		mysql2.GenMySQL(c)
	case base.CodeDBMSSQL:

	default:
		log.Panicln("database type err")
	}
}
