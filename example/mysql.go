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
			// 1：不生成数据库基本信息 2：生成简单的数据库字段信息
			GenDBInfoType: 1,
			// json tag类型，前提：IsGenJsonTag:true. 1.UserName 2.user_name 3.userName 4.user-name
			JsonTagType: 3,
			// 是否生成建表语句
			IsGenCreateSQL: true})
	time.Sleep(time.Second)
}
