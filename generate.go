package ormtool

import (
	"github.com/enjoy322/ormtool/base"
	"github.com/enjoy322/ormtool/mysqlTool"
	"log"
)

func GenerateMySQL(dbConf base.MysqlConfig, c base.Config) {
	log.Println("-----generate-----")
	mysqlTool.GenMySQL(dbConf, c)
	log.Println("-----finish-----")
}
