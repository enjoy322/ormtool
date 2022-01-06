package ormtool

import (
	"convert/base"
	"convert/mysqlTool"
	"log"
)

func GenerateMySQL(dbConf base.MysqlConfig, c base.Config) {
	log.Println("-----generate-----")
	mysqlTool.GenMySQL(dbConf, c)
	log.Println("-----finish-----")
}
