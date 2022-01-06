package mysqlTool

import (
	"database/sql"
	"log"
	"ormtool/base"
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
	IsNullable    string
	ColumnKey     string
	Tag           string
}

func (s service) StructContent(dbName string, c base.Config) (packageName, fileDir, fileName string, data map[string]string) {
	tables := s.DealColumn(c)
	packageName, fileDir, fileName = base.DealFilePath(c.SavePath, dbName)
	data = make(map[string]string)
	for tableName, columns := range tables {
		var structInfo strings.Builder
		structName := tableName
		if len(structName) == 1 {
			structName = strings.ToUpper(tableName[0:1])
		} else {
			split := strings.Split(tableName, "_")
			var tName strings.Builder
			for _, str := range split {
				tName.WriteString(strings.ToUpper(str[0:1]) + str[1:])
			}
			structName = tName.String()
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
		data[tableName] = structInfo.String()
	}
	return
}

func (s service) DealColumn(c base.Config) map[string][]column {
	tables := s.GetColumn()
	for _, cols := range tables {
		for i, col := range cols {
			var f bool
			if c.IsGenJsonTag {
				f = true
				cols[i].Tag = col.ColumnName
				cols[i].Tag = strings.ToLower(cols[i].Tag)
				cols[i].Tag = "`json:\"" + cols[i].Tag + "\" "
			}
			switch c.GenDBInfoType {
			case base.CodeDBInfoNone:

			case base.CodeDBInfoSimple:
				if !f {
					cols[i].Tag += "`"
				}
				f = true
				cols[i].Tag += "db:\"" + col.ColumnType
				var sNull string
				if col.IsNullable == "NO" {
					sNull = " not null"
				}
				cols[i].Tag += sNull
				if col.Default != nil {
					cols[i].Tag += " default " + string(col.Default.([]uint8))
				}
				cols[i].Tag += "\""

			case base.CodeDBInfoGorm:
				if !f {
					cols[i].Tag += "`"
				}
				f = true
				cols[i].Tag += "gorm:\"column:" + col.ColumnName
				cols[i].Tag += ";type:" + col.ColumnType
				if col.Default != nil {
					cols[i].Tag += " default:" + string(col.Default.([]uint8))
				}
				var sNull string
				if col.IsNullable == "NO" {
					sNull = " not null"
				}
				cols[i].Tag += sNull
				cols[i].Tag += "\""

			case base.CodeDBInfoXorm:

			}

			if f {
				cols[i].Tag += "`"
			}
			cols[i].ColumnName = base.CamelCase(col.ColumnName)
			if col.ColumnType == "tinyint(1)" {
				cols[i].ColumnType = "bool"
			} else {
				cols[i].ColumnType = mysqlToGo[col.DataType]
			}
		}
	}
	return tables
}

// GetColumn 获取数据库表信息
func (s service) GetColumn() map[string][]column {
	/*
		# SELECT * FROM information_schema.COLUMNS WHERE table_schema = DATABASE();
		# SELECT TABLE_NAME,TABLE_COMMENT FROM information_schema.TABLES WHERE table_schema = DATABASE();

		# SELECT * FROM information_schema.COLUMNS WHERE table_schema = DATABASE();

		show indexes from user_card
	*/
	//todo 表约束信息
	tables := make(map[string][]column)
	//IS_NULLABLE,COLUMN_DEFAULT,CHARACTER_MAXIMUM_LENGTH
	sqlStr := "SELECT COLUMN_NAME,DATA_TYPE,COLUMN_TYPE,COLUMN_DEFAULT,TABLE_NAME," +
		"COLUMN_COMMENT,character_maximum_length,IS_NULLABLE,COLUMN_KEY" +
		" FROM information_schema.COLUMNS WHERE table_schema = DATABASE()"
	rows, err := s.DB.Query(sqlStr)
	if err != nil {
		log.Panicln(err.Error())
	}
	defer func(rows *sql.Rows) {
		err := rows.Close()
		if err != nil {

		}
	}(rows)
	for rows.Next() {
		col := column{}
		err = rows.Scan(&col.ColumnName, &col.DataType, &col.ColumnType, &col.Default,
			&col.TableName, &col.ColumnComment, &col.Length, &col.IsNullable, &col.ColumnKey)
		if err != nil {
			log.Println(err.Error())
		}
		tables[col.TableName] = append(tables[col.TableName], col)
	}
	return tables
}