package ormtool

import (
	"log"

	"github.com/enjoy322/ormtool/base"
	"github.com/enjoy322/ormtool/mysqlTool"
)

func GenerateMySQL(dbConf base.MysqlConfig, c base.Config) {
	log.Println("-----ormtool generating-----")
	mysqlTool.GenMySQL(dbConf, c)
	log.Println("-----ormtool done-----")
}
