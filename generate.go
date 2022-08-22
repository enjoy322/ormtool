package ormtool

import (
	"log"

	"github.com/enjoy322/ormtool/base"
	"github.com/enjoy322/ormtool/mysqlTool"
)

func GenerateMySQL(c base.Config) {
	log.Println("-----ormtool generating-----")
	mysqlTool.NewService(c).Gen()
	log.Println("-----ormtool done-----")
}
