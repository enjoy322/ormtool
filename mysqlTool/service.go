package mysqlTool

import (
	"database/sql"
	"log"
	"strings"

	"github.com/enjoy322/ormtool/base"
)

type Method interface {
	GenStruct() (fileData base.FileInfo, data []base.StructInfo)

	listTables() []tableInfo

	listColumns() map[string][]column

	dealColumn(t *tableInfo)

	dealStructContent(t tableInfo) string

	dealType(c base.Config, typeSimple, typeDetail string) string

	getCreateSQL(tableName string) string
}
type service struct {
	DB       *sql.DB
	Info     []base.StructInfo
	FileInfo base.FileInfo
	dbName   string
	Conf     base.Config
}

func NewService(c base.Config) *service {
	conn := dbConn(c.ConnStr)

	return &service{DB: conn, dbName: c.Database, Conf: c}
}

func (s service) Gen() {
	defer func(db *sql.DB) {
		err := db.Close()
		if err != nil {
			log.Fatalln(err)
		}
	}(s.DB)

	fileData, data := s.genStruct()
	// write into file
	base.Write(fileData, data, s.Conf.IsGenInOneFile)

}

type tableInfo struct {
	TableName    string
	TableComment string
	column       []column
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
	//Length        interface{}
	IsNullable string
	//ColumnKey     string
	Tag string
}

// GenStruct struct info, include: struct comment, create table sql
func (s service) genStruct() (fileData base.FileInfo, data []base.StructInfo) {
	// save file info
	s.FileInfo.PackageName, s.FileInfo.FileDir, s.FileInfo.FileName = base.DealFilePath(s.Conf.SavePath, s.dbName)

	tables := s.listTables()
	columns := s.listColumns()

	for _, table := range tables {
		if v, ok := columns[table.TableName]; ok {
			table.column = v
		}

		var info base.StructInfo

		// table name
		info.TableName = table.TableName

		// create table sql
		if s.Conf.IsGenCreateSQL {
			info.CreateSQL = s.getCreateSQL(table.TableName)
		}
		// deal column
		s.dealColumn(&table)

		// struct info
		info.StructContent = s.dealStructContent(table)

		info.Name = base.UpperCamel(table.TableName)

		// table comment
		// add if table comment exists
		if table.TableComment != "" || s.Conf.IsGenCreateSQL {
			info.Note = "// " + info.Name + "\t" + table.TableComment + "\n"
		}

		s.Info = append(s.Info, info)

	}
	return s.FileInfo, s.Info
}

// DealColumn judge column type and generate tag info
func (s service) dealColumn(t *tableInfo) {
	for i := 0; i < len(t.column); i++ {
		var f bool
		if s.Conf.IsGenJsonTag {
			//生成 json tag
			f = true
			t.column[i].Tag = "`json:\"" + base.JsonTag(s.Conf.JsonTagType, t.column[i].ColumnName) + "\""
		}
		if s.Conf.GenDBInfoType == 2 {
			t.column[i].Tag = t.column[i].Tag + " "
		}
		switch s.Conf.GenDBInfoType {
		case 1:
		case 2:
			if !f {
				t.column[i].Tag += "`"
			}
			f = true
			t.column[i].Tag += "db:\"" + t.column[i].ColumnType
			var sNull string
			if t.column[i].IsNullable == "NO" {
				sNull = " not null"
			}
			t.column[i].Tag += sNull
			if t.column[i].Default != nil {
				t.column[i].Tag += " default " + string(t.column[i].Default.([]uint8))
			}
			t.column[i].Tag += "\""
		}

		if f {
			t.column[i].Tag += "`"
		}
		t.column[i].ColumnName = base.UpperCamel(t.column[i].ColumnName)
		t.column[i].ColumnType = s.dealType(s.Conf, t.column[i].DataType, t.column[i].ColumnType)
	}
}

func (s service) dealStructContent(t tableInfo) string {
	var info strings.Builder
	// struct name
	structName := base.UpperCamel(t.TableName)

	info.WriteString("type " + structName + " struct {\n")
	for _, v := range t.column {
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
	info.WriteString("return \"" + t.TableName + "\"")
	info.WriteString("\n}\n")

	info.WriteString("var " + structName + "Col = struct {\n")
	for _, v := range t.column {
		info.WriteString(v.ColumnName)
		info.WriteString("\t" + "string\n")
	}
	info.WriteString("}{\n")
	for _, v := range t.column {
		info.WriteString(v.ColumnName)
		info.WriteString(":\t\"" + strings.ToLower(v.ColumnDBName) + "\"" + ",\n")
	}
	info.WriteString("\n}\n")

	return info.String()
}

func (s service) dealType(c base.Config, typeSimple, typeDetail string) string {
	if v, ok := c.CustomType[typeDetail]; ok {
		return v
	}
	switch typeSimple {
	case "int":
		return mysqlToGo[typeDetail]
	default:
		return mysqlToGo[typeSimple]
	}
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

func (s service) listColumns() map[string][]column {
	tables := make(map[string][]column)
	sqlStr := `SELECT COLUMN_NAME,DATA_TYPE,COLUMN_TYPE,COLUMN_DEFAULT,TABLE_NAME,
       COLUMN_COMMENT
   FROM information_schema.COLUMNS WHERE table_schema = ? order by TABLE_NAME`
	rows, err := s.DB.Query(sqlStr, s.dbName)

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
			&col.TableName, &col.ColumnComment)
		if err != nil {
			log.Fatalln(err.Error())
		}
		col.ColumnDBName = col.ColumnName
		tables[col.TableName] = append(tables[col.TableName], col)
	}
	return tables
}

func (s service) listTables() []tableInfo {
	sqlStr := `select Table_Name,Table_Comment from information_schema.TABLES where TABLE_SCHEMA=?`

	rows, err := s.DB.Query(sqlStr, s.dbName)
	if err != nil {
		log.Fatalln(err.Error())
	}
	defer func(rows *sql.Rows) {
		err := rows.Close()
		if err != nil {
			log.Fatalln(err)
		}
	}(rows)
	var list []tableInfo
	for rows.Next() {
		var i tableInfo
		err := rows.Scan(&i.TableName, &i.TableComment)
		if err != nil {
			log.Fatalln(err)
		}
		list = append(list, i)
	}
	return list
}
