package ormtool

import (
	"testing"
)

func TestGenerateMySQL(t *testing.T) {
	GenerateMySQL(
		Config{
			//[user]:[password]@tcp([host]:[port])/[database]?parseTime=true
			ConnStr: "root:qwe123@tcp(127.0.0.1:3306)/orm-test?parseTime=true",
			// database name
			Database: "orm-test",
			// relative path
			SavePath: "./model/model.go",
			// Generate one file or files
			IsGenInOneFile: false,
			// Generate simple database field information like: "int unsigned not null"
			// value 1:not generate; 2：simple info
			GenDBInfoType: 1,
			// json tag
			IsGenJsonTag: true,
			// json tag type. The necessary conditions：IsGenJsonTag:true.
			// 1.UserName 2.userName 3.user_name 4.user-name
			JsonTagType: 3,
			// sql of creating table in database
			IsGenCreateSQL: false,
			// simple crud function
			IsGenFunction: true,
			// cache simple model info to redis, Ps. IsGenCreateSQL = true
			IsGenFunctionWithCache: true,
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
