package example

import (
	"github.com/enjoy322/ormtool"
	"github.com/enjoy322/ormtool/base"
)

func GenMysql() {
	ormtool.GenerateMySQL(
		base.MysqlConfig{
			User:     "root",
			Password: "qwe123",
			Host:     "127.0.0.1",
			Port:     "3306",
			Database: "demo"},
		base.Config{
			// relative path
			SavePath: "../models/model.go",
			// Generate one file or files
			IsGenInOneFile: true,
			// Generate simple database field information like: "int unsigned not null"
			// value 1:not generate; 2：simple info
			GenDBInfoType: 1,
			// json tag
			IsGenJsonTag: true,
			// json tag type. The necessary conditions：IsGenJsonTag:true.
			// 1.UserName 2.userName 3.user_name 4.user-name
			JsonTagType: 3,
			// sql of creating table in database
			IsGenCreateSQL: true,
			// custom type relationships will be preferred
			// the key is the database type, the value is the golang type
			CustomType: map[string]string{
				"int":          "int",
				"int unsigned": "uint32",
				"tinyint(1)":   "bool",
				"json":         "json.RawMessage",
			},
		})
}
