package main

import (
	"ormtool"
	"ormtool/base"
	"time"
)

func main() {
	ormtool.GenerateMySQL(
		base.MysqlConfig{
			User:     "root",
			Password: "qwe123",
			Host:     "127.0.0.1",
			Port:     "3306",
			Database: "test"},
		base.Config{
			SavePath:       "./models/test.go",
			IsGenJsonTag:   true,
			IsGenInOneFile: true,
			GenDBInfoType:  base.CodeDBInfoSimple,
			JsonTagType:    base.CodeJsonTag1})
	time.Sleep(time.Second)
}
