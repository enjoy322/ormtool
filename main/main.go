package main

import (
	"convert"
	"convert/base"
	"time"
)

func main() {
	convert.Gen(base.Config{
		DataBaseType: base.CodeDBMySQL,
		MySQL: base.MysqlConfig{
			User:     "root",
			Password: "qwe123",
			Host:     "127.0.0.1",
			Port:     "3306",
			Database: "test"},
		SavePath:       "./models/test.go",
		IsGenJsonTag:   true,
		IsGenInOneFile: true,
		GenDBInfoType:  1,
		JsonTagType:    base.CodeJsonTag1})
	time.Sleep(time.Second)
}
