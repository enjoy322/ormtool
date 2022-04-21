package mysqlTool

import (
	"database/sql"
	"github.com/enjoy322/ormtool/base"
	"log"
	"strconv"
	"strings"
)

type service struct {
	*sql.DB
}

func Service(DB *sql.DB) *service {
	return &service{DB: DB}
}

type column struct {
	ColumnDBName string
	//列名
	ColumnName string
	//字段类型，如varchar
	DataType string
	//字段类型，显示细节，如varchar(32)
	ColumnType string
	//默认值
	Default interface{}
	//表名
	TableName string
	//字段注释
	ColumnComment string
	Length        interface{}
	IsNullable    string
	ColumnKey     string
	Tag           string
}

// StructContent 结构体信息
func (s service) StructContent(dbName string, c base.Config) (packageName, fileDir, fileName string, data map[string]string) {
	// 查询表名注释
	tableCommentMap := s.GetTableComment(dbName)
	// 查询数据库的所有表
	tables := s.DealColumn(c)
	packageName, fileDir, fileName = base.DealFilePath(c.SavePath, dbName)
	data = make(map[string]string)
	for tableName, columns := range tables {
		var createSQL string
		if c.IsGenCreateSQL {
			// 需要建表SQL语句
			createSQL = s.GetCreateSQL(tableName)
		}

		// 表（结构体内容）
		var structInfo strings.Builder
		// 结构体名称
		structName := tableName
		if len(structName) == 1 {
			structName = strings.ToUpper(tableName[:1])
		} else {
			split := strings.Split(tableName, "_")
			var tName strings.Builder
			for _, str := range split {
				tName.WriteString(strings.ToUpper(str[:1]) + str[1:])
			}
			structName = tName.String()
		}

		// 结构体名称后加注释（如果表存在注释情况下
		if v, ok := tableCommentMap[tableName]; ok {
			if v != "" || c.IsGenCreateSQL {
				//判断生成表注释
				structInfo.WriteString("// " + structName + "\t" + v + "\n")
			}
		}

		//添加建表SQL语句
		if c.IsGenCreateSQL {
			structInfo.WriteString("/*")
			structInfo.WriteString(createSQL)
			structInfo.WriteString("*/\n")
		}

		// 结构体字段
		structInfo.WriteString("type " + structName + " struct {\n")
		for _, v := range columns {
			structInfo.WriteString("\t")
			structInfo.WriteString(v.ColumnName)
			structInfo.WriteString("\t")
			structInfo.WriteString(v.ColumnType)
			structInfo.WriteString("\t")
			structInfo.WriteString(v.Tag)
			structInfo.WriteString("\t")
			if v.ColumnComment != "" {
				structInfo.WriteString(" // ")
				structInfo.WriteString(v.ColumnComment)
			}
			structInfo.WriteString("\n")
		}
		structInfo.WriteString("}\n\n")
		// 数据库表名函数
		structInfo.WriteString("func (*" + structName + ") TableName() string {\n")
		structInfo.WriteString("return \"" + tableName + "\"")
		structInfo.WriteString("\n}\n")

		//结构体字段与表字段对应
		structInfo.WriteString("var " + structName + "Col = struct {\n")
		for _, v := range columns {
			structInfo.WriteString(v.ColumnName)
			structInfo.WriteString("\t" + "string\n")
		}
		structInfo.WriteString("}{\n")
		for _, v := range columns {
			structInfo.WriteString(v.ColumnName)
			structInfo.WriteString(":\t\"" + strings.ToLower(v.ColumnDBName) + "\"" + ",\n")
		}
		structInfo.WriteString("\n}\n")

		data[tableName] = structInfo.String()
	}
	return
}

// DealColumn 处理结构体字段 生成的tag信息
func (s service) DealColumn(c base.Config) map[string][]column {
	tables := s.GetColumn()
	for _, cols := range tables {
		for i, col := range cols {
			var f bool
			if c.IsGenJsonTag {
				//生成 json tag
				f = true
				cols[i].Tag = "`json:\"" + jsonTag(c.JsonTagType, col.ColumnName) + "\" "
			}
			switch c.GenDBInfoType {
			case 1:

			case 2:
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
			}

			if f {
				cols[i].Tag += "`"
			}
			cols[i].ColumnName = base.CamelCase(col.ColumnName)
			cols[i].ColumnType = dealType(c, col.DataType, col.ColumnType)
		}
	}
	return tables
}

func dealType(c base.Config, typeSimple, typeDetail string) string {
	if v, ok := c.CustomType[typeDetail]; ok {
		return v
	}
	switch typeSimple {
	case "tinyint":
		num := getTypeNum(typeDetail)
		switch num {
		case 0:
			return "string"
		case 1:
			return "bool"
		}
	case "int":
		return mysqlToGo[typeDetail]
	default:
		return mysqlToGo[typeSimple]
	}
	return ""
}

// 获取表字段长度约束
func getTypeNum(typeStr string) int {
	f := strings.HasSuffix(typeStr, ")")
	if f {
		//	有长度约束
		splitAfter := strings.SplitAfter(typeStr, "(")
		n := splitAfter[1][:1]
		i, err := strconv.Atoi(n)
		if err != nil {
			panic(err)
		}
		return i
	}
	return 0
}

// 处理tag： json
func jsonTag(jsonType int, origin string) string {
	switch jsonType {
	//1.UserName 2.userName 3.user_name 4.user-name
	case 1:
		return Case2Camel(origin)
	case 2:
		s1 := Case2Camel(origin)
		return strings.ToLower(s1[:1]) + s1[1:]
	case 3:
		return strings.ToLower(origin)
	case 4:
		return strings.Replace(origin, "_", "-", -1)

	}
	panic("json tag 参数错误")
}

func Case2Camel(name string) string {
	name = strings.Replace(name, "_", " ", -1)
	name = strings.Title(name)
	return strings.Replace(name, " ", "", -1)
}

// GetColumn 获取数据库表信息
func (s service) GetColumn() map[string][]column {
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
		col.ColumnDBName = col.ColumnName
		tables[col.TableName] = append(tables[col.TableName], col)
	}
	return tables
}

type CreateSQL struct {
	Table string `json:"Table"`
	SQL   string `json:"Create Table"`
}

// GetCreateSQL 获取建表语句
func (s service) GetCreateSQL(tableName string) string {
	sqlStr := "show create table " + tableName
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
		var cSql CreateSQL
		err = rows.Scan(&cSql.Table, &cSql.SQL)
		if err != nil {
			log.Println(err.Error())
		}
		return cSql.SQL
	}
	return ""
}

// GetTableComment  获取表信息
func (s service) GetTableComment(dbName string) map[string]string {
	sqlStr := "show table status from " + dbName
	rows, err := s.DB.Query(sqlStr)
	if err != nil {
		log.Panicln(err.Error())
	}
	defer func(rows *sql.Rows) {
		err := rows.Close()
		if err != nil {

		}
	}(rows)
	columns, _ := rows.Columns()
	columnLength := len(columns)
	cache := make([]interface{}, columnLength) //临时存储每行数据
	for index, _ := range cache {              //为每一列初始化一个指针
		var a interface{}
		cache[index] = &a
	}

	var list []map[string]interface{} //返回的切片
	for rows.Next() {
		err = rows.Scan(cache...)
		if err != nil {
			log.Println(err.Error())
		}
		item := make(map[string]interface{})
		for i, data := range cache {
			item[columns[i]] = *data.(*interface{}) //取实际类型
		}
		list = append(list, item)
	}
	m := make(map[string]string)
	for _, i := range list {
		if v, ok := i["Name"]; ok {
			tName := string(v.([]uint8))
			if v, ok := i["Comment"]; ok {
				comment := string(v.([]uint8))
				m[tName] = comment
			}
		}
	}
	return m
}
