package ormtool

import (
	"log"
	"ormtool/base"
	"ormtool/mysqlTool"
)

func GenerateMySQL(dbConf base.MysqlConfig, c base.Config) {
	log.Println("-----generate-----")
	mysqlTool.GenMySQL(dbConf, c)
	log.Println("-----finish-----")
}
