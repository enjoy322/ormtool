package mysql2

import (
	"convert/base"
	"database/sql"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"
)

type service struct {
	*sql.DB
}

func Service(DB *sql.DB) *service {
	return &service{DB: DB}
}

type column struct {
	ColumnName    string
	DataType      string
	ColumnType    string
	Default       interface{}
	TableName     string
	ColumnComment string
	Length        interface{}
	Tag           string
}

func dealFilePath(s string, db string) (packageName, fileDir, fileName string) {
	if !strings.HasSuffix(s, ".go") {
		fmt.Println("保存路径错误，正确如./models/xx.go")
		os.Exit(0)
	}
	if len(strings.Trim(s, " ")) < 1 {
		packageName = "models"
		fileDir = "models"
		fileName = db
		return
	}
	split := strings.Split(s, "/")
	if len(split) <= 1 {
		packageName = "models"
		fileDir = "models"
		fileName = s
	} else {
		packageName = split[len(split)-2]
		fileName = split[len(split)-1]
		s2 := strings.Split(s, "/"+fileName)
		fileDir = s2[0]
	}
	return
}

func (s service) Create(c base.Config) {
	tables := s.GetColumn(c)
	packageName, fileDir, _ := dealFilePath(c.SavePath, c.MySQL.Database)
	packageStr := "package " + packageName + "\n\n"
	os.MkdirAll(fileDir, 0777)
	for tableName, columns := range tables {
		var structInfo strings.Builder
		structInfo.WriteString(packageStr)
		structName := tableName
		if len(structName) == 1 {
			structName = strings.ToUpper(tableName[0:1])
		} else {
			structName = strings.ToUpper(tableName[0:1]) + tableName[1:]
		}
		depth := 1
		structInfo.WriteString("type " + structName + " struct {\n")
		for _, v := range columns {
			structInfo.WriteString(base.Tab(depth))
			structInfo.WriteString(v.ColumnName)
			structInfo.WriteString(base.Tab(depth))
			structInfo.WriteString(v.ColumnType)
			structInfo.WriteString(base.Tab(depth))
			structInfo.WriteString(v.Tag)
			structInfo.WriteString(base.Tab(depth))
			if v.ColumnComment != "" {
				structInfo.WriteString(" // ")
				structInfo.WriteString(v.ColumnComment)
			}
			structInfo.WriteString(base.Next(1))
		}
		structInfo.WriteString(base.Tab(depth-1) + "}\n\n")

		// 数据库表名函数
		structInfo.WriteString("func (" + structName + ") TableName() string {\n")
		structInfo.WriteString("return \"" + tableName + "\"")
		structInfo.WriteString("\n}\n")
		fileName := fileDir + "/" + tableName + ".go"
		f, err := os.Create(fileName)
		if err != nil {
			log.Panicln(err.Error())
		}
		defer f.Close()
		f.WriteString(structInfo.String())

		cmd := exec.Command("gofmt", "-w", fileName)
		cmd.Run()
	}

	log.Println("finish!")
}

func (s service) CreateOneFile(c base.Config) {
	tables := s.GetColumn(c)
	packageName, fileDir, fileName := dealFilePath(c.SavePath, c.MySQL.Database)
	var structInfo strings.Builder
	structInfo.WriteString("package " + packageName + "\n\n")
	for tableName, columns := range tables {
		structName := tableName
		if len(structName) == 1 {
			structName = strings.ToUpper(tableName[0:1])
		} else {
			structName = strings.ToUpper(tableName[0:1]) + tableName[1:]
		}
		depth := 1
		structInfo.WriteString("type " + structName + " struct {\n")
		for _, v := range columns {
			structInfo.WriteString(base.Tab(depth))
			structInfo.WriteString(v.ColumnName)
			structInfo.WriteString(base.Tab(depth))
			structInfo.WriteString(v.ColumnType)
			structInfo.WriteString(base.Tab(depth))
			structInfo.WriteString(v.Tag)
			structInfo.WriteString(base.Tab(depth))
			if v.ColumnComment != "" {
				structInfo.WriteString(" // ")
				structInfo.WriteString(v.ColumnComment)
			}
			structInfo.WriteString(base.Next(1))
		}
		structInfo.WriteString(base.Tab(depth-1) + "}\n\n")

		// 数据库表名函数
		structInfo.WriteString("func (" + structName + ") TableName() string {\n")
		structInfo.WriteString("return \"" + tableName + "\"")
		structInfo.WriteString("\n}\n")
	}
	os.MkdirAll(fileDir, 0777)
	f, err := os.Create(fileDir + "/" + fileName)
	if err != nil {
		log.Panicln(err.Error())
	}
	defer f.Close()
	f.WriteString(structInfo.String())

	cmd := exec.Command("gofmt", "-w", fileDir+"/"+fileName)
	cmd.Run()
	log.Println("finish!")
}
func (s service) GetColumn(c base.Config) map[string][]column {
	columns := make(map[string][]column)
	//IS_NULLABLE,COLUMN_DEFAULT,CHARACTER_MAXIMUM_LENGTH
	sqlStr := "SELECT COLUMN_NAME,DATA_TYPE,COLUMN_TYPE,COLUMN_DEFAULT,TABLE_NAME,COLUMN_COMMENT,character_maximum_length" +
		" FROM information_schema.COLUMNS WHERE table_schema = DATABASE()"
	rows, err := s.DB.Query(sqlStr)
	if err != nil {
		log.Panicln(err.Error())
	}
	defer rows.Close()
	for rows.Next() {
		col := column{}
		err = rows.Scan(&col.ColumnName, &col.DataType, &col.ColumnType, &col.Default, &col.TableName, &col.ColumnComment, &col.Length)
		if err != nil {
			log.Println(err.Error())
		}
		var f bool
		if c.IsGenJsonTag {
			f = true
			col.Tag = col.ColumnName
			col.Tag = strings.ToLower(col.Tag)
			col.Tag = "`json:\"" + col.Tag + "\" "
		}
		switch c.GenDBInfoType {
		case base.CodeDBInfoSimple:
			if !f {
				col.Tag += "`"
			}
			f = true
			col.Tag += "db:\"" + col.ColumnType + "\""

		}
		if f {
			col.Tag += "`"
		}

		col.ColumnName = base.CamelCase(col.ColumnName)
		if strings.HasPrefix(col.ColumnType, "varchar") {
			col.ColumnType = "string"
		} else {
			col.ColumnType = MysqlToGo[col.ColumnType]
		}
		columns[col.TableName] = append(columns[col.TableName], col)
	}
	return columns
}
