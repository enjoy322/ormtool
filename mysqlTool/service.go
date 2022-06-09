package mysqlTool

import (
	"database/sql"
	"log"
	"strings"

	"github.com/enjoy322/ormtool/base"
)

type Method interface {
	GenStruct(dbName string, c base.Config) (fileData base.FileInfo, data []base.StructInfo)

	dealStructContent(tableName string, columns []column) string

	dealColumn(c base.Config) map[string][]column

	dealType(c base.Config, typeSimple, typeDetail string) string

	getColumn() map[string][]column

	getTableComment(dbName string) map[string]string

	getCreateSQL(tableName string) string
}
type service struct {
	DB       *sql.DB
	Info     []base.StructInfo
	FileInfo base.FileInfo
}

func Service(DB *sql.DB) *service {
	return &service{DB: DB}
}

type column struct {
	ColumnDBName string
	ColumnName   string
	//example: varchar
	DataType string
	//example: varchar(32)
	ColumnType string
	//default value
	Default       interface{}
	TableName     string
	ColumnComment string
	Length        interface{}
	IsNullable    string
	ColumnKey     string
	Tag           string
}

// GenStruct struct info, include: struct comment, create table sql
func (s service) GenStruct(dbName string, c base.Config) (fileData base.FileInfo, data []base.StructInfo) {
	// all table comments
	tableCommentMap := s.getTableComment(dbName)
	// all tables
	tables := s.dealColumn(c)
	// save file info
	s.FileInfo.PackageName, s.FileInfo.FileDir, s.FileInfo.FileName = DealFilePath(c.SavePath, dbName)

	// data = make(map[string]string)

	for tableName, columns := range tables {
		var info base.StructInfo

		// table name
		info.TableName = tableName

		// create table sql
		if c.IsGenCreateSQL {
			info.CreateSQL = s.getCreateSQL(tableName)
		}

		// struct info
		info.StructContent = s.dealStructContent(tableName, columns)

		info.Name = base.UpperCamel(tableName)

		// table comment
		// add if table comment exists
		if v, ok := tableCommentMap[tableName]; ok {
			if v != "" || c.IsGenCreateSQL {
				info.Note = ("// " + info.Name + "\t" + v + "\n")
			}
		}

		s.Info = append(s.Info, info)

	}
	return s.FileInfo, s.Info
}

func (s service) dealStructContent(tableName string, columns []column) string {
	var info strings.Builder
	// struct name
	structName := base.UpperCamel(tableName)

	info.WriteString("type " + structName + " struct {\n")
	for _, v := range columns {
		info.WriteString("\t")
		info.WriteString(v.ColumnName)
		info.WriteString("\t")
		info.WriteString(v.ColumnType)
		info.WriteString("\t")
		info.WriteString(v.Tag)
		info.WriteString("\t")
		if v.ColumnComment != "" {
			info.WriteString(" // ")
			info.WriteString(v.ColumnComment)
		}
		info.WriteString("\n")
	}
	info.WriteString("}\n\n")
	// function for get table name in database
	info.WriteString("func (*" + structName + ") TableName() string {\n")
	info.WriteString("return \"" + tableName + "\"")
	info.WriteString("\n}\n")

	info.WriteString("var " + structName + "Col = struct {\n")
	for _, v := range columns {
		info.WriteString(v.ColumnName)
		info.WriteString("\t" + "string\n")
	}
	info.WriteString("}{\n")
	for _, v := range columns {
		info.WriteString(v.ColumnName)
		info.WriteString(":\t\"" + strings.ToLower(v.ColumnDBName) + "\"" + ",\n")
	}
	info.WriteString("\n}\n")

	return info.String()
}

// DealColumn judge column type and generate tag info
func (s service) dealColumn(c base.Config) map[string][]column {
	tables := s.getColumn()
	for _, cols := range tables {
		for i, col := range cols {
			var f bool
			if c.IsGenJsonTag {
				//生成 json tag
				f = true
				cols[i].Tag = "`json:\"" + base.JsonTag(c.JsonTagType, col.ColumnName) + "\" "
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
			cols[i].ColumnName = base.UpperCamel(col.ColumnName)
			cols[i].ColumnType = s.dealType(c, col.DataType, col.ColumnType)
		}
	}
	return tables
}

func (s service) dealType(c base.Config, typeSimple, typeDetail string) string {
	if v, ok := c.CustomType[typeDetail]; ok {
		return v
	}
	switch typeSimple {
	case "tinyint":
		num := base.GetTypeNum(typeDetail)
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

// GetColumn columns of table
func (s service) getColumn() map[string][]column {
	tables := make(map[string][]column)
	//IS_NULLABLE,COLUMN_DEFAULT,CHARACTER_MAXIMUM_LENGTH
	sqlStr := "SELECT COLUMN_NAME,DATA_TYPE,COLUMN_TYPE,COLUMN_DEFAULT,TABLE_NAME," +
		"COLUMN_COMMENT,character_maximum_length,IS_NULLABLE,COLUMN_KEY" +
		" FROM information_schema.COLUMNS WHERE table_schema = DATABASE()"
	rows, err := s.DB.Query(sqlStr)
	if err != nil {
		log.Fatalln(err.Error())
	}
	defer func(rows *sql.Rows) {
		err := rows.Close()
		if err != nil {
			log.Fatalln(err)
		}
	}(rows)
	for rows.Next() {
		col := column{}
		err = rows.Scan(&col.ColumnName, &col.DataType, &col.ColumnType, &col.Default,
			&col.TableName, &col.ColumnComment, &col.Length, &col.IsNullable, &col.ColumnKey)
		if err != nil {
			log.Fatalln(err.Error())
		}
		col.ColumnDBName = col.ColumnName
		tables[col.TableName] = append(tables[col.TableName], col)
	}
	return tables
}

// GetCreateSQL sql of creating table in database
func (s service) getCreateSQL(tableName string) string {
	sqlStr := "show create table " + tableName
	rows, err := s.DB.Query(sqlStr)
	if err != nil {
		log.Fatalln(err.Error())
	}
	defer func(rows *sql.Rows) {
		err := rows.Close()
		if err != nil {
			log.Fatalln(err)
		}
	}(rows)

	type CreateSQL struct {
		Table string `json:"Table"`
		SQL   string `json:"Create Table"`
	}
	var cSql CreateSQL

	for rows.Next() {
		err = rows.Scan(&cSql.Table, &cSql.SQL)
		if err != nil {
			log.Fatalln(err.Error())
		}
	}

	var info strings.Builder
	info.WriteString("/*")
	info.WriteString(cSql.SQL)
	info.WriteString("*/\n")

	return info.String()
}

// GetTableComment  comment fo table
func (s service) getTableComment(dbName string) map[string]string {
	sqlStr := "show table status from " + dbName
	rows, err := s.DB.Query(sqlStr)
	if err != nil {
		log.Fatalln(err.Error())
	}
	defer func(rows *sql.Rows) {
		err := rows.Close()
		if err != nil {
			log.Fatalln(err)
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
			log.Fatalln(err.Error())
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

// DealFilePath back save path and package name
func DealFilePath(s string, db string) (packageName, fileDir, fileName string) {
	if !strings.HasSuffix(s, ".go") {
		log.Fatalln("path error! correct example: ./models/xx.go")
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
