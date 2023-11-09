package ormtool

import (
	"github.com/enjoy322/ormtool/base"
	"github.com/enjoy322/ormtool/mysqlTool"
)

func GenerateMySQL(c base.Config) {
	mysqlTool.NewService(c).Gen()
}
