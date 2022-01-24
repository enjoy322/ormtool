package main

import (
	"github.com/enjoy322/ormtool"
	"github.com/enjoy322/ormtool/base"
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
			SavePath:       "./models/model.go",
			IsGenJsonTag:   true,
			IsGenInOneFile: true,
			GenDBInfoType:  base.CodeDBInfoNone,
			JsonTagType:    base.CodeJsonTag1,
			IsGenCreateSQL: false})
	time.Sleep(time.Second)
}
