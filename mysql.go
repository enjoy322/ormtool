package ormtool

import (
	"database/sql"
	"fmt"
	"log"
	"strings"
)

type Method interface {
	genStruct() (fileData FileInfo, data []StructInfo)

	listTables() []tableInfo

	listColumns() map[string][]column

	dealColumn(t *tableInfo)

	dealStructContent(t tableInfo) string

	dealType(c Config, typeSimple, typeDetail string) string

	getCreateSQL(tableName string) string
}

func GenerateMySQL(c Config) {
	log.Println("-----generate start-----")
	newService(c).Gen()
	log.Println("-----generate done-----")
}

type service struct {
	DB   *sql.DB
	Info []StructInfo
	// save model file
	FileSave FileInfo
	dbName   string
	Conf     Config
}

func newService(c Config) *service {
	conn := mysqlConn(c.ConnStr)

	return &service{DB: conn, dbName: c.Database, Conf: c}
}

func (s service) Gen() {
	defer func(db *sql.DB) {
		err := db.Close()
		if err != nil {
			log.Fatalln(err)
		}
	}(s.DB)

	fileSave, data := s.genStruct()
	// write into file
	Write(fileSave, data, s.Conf.IsGenInOneFile)

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
func (s service) genStruct() (fileSave FileInfo, data []StructInfo) {
	// save file info
	s.FileSave.PackageName, s.FileSave.FileDir, s.FileSave.FileName = DealFilePath(s.Conf.SavePath, s.dbName)

	tables := s.listTables()
	columns := s.listColumns()

	// deal per table
	for _, table := range tables {
		if v, ok := columns[table.TableName]; ok {
			table.column = v
		}

		var info StructInfo

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

		info.Name = UpperCamel(table.TableName)

		// table comment
		// add if table comment exists
		if table.TableComment != "" || s.Conf.IsGenCreateSQL {
			info.Note = "// " + info.Name + "\t" + table.TableComment + "\n"
		}

		if s.Conf.IsGenFunction {
			//	simple function
			info.Function = s.genFunction(info.Name)
		}

		s.Info = append(s.Info, info)
	}
	return s.FileSave, s.Info
}
func (s service) genFunction(name string) string {

	var info strings.Builder
	info.WriteString("// function\n")
	interfaceName := name + "ModelInterface"
	interfaceContent := `
type %s interface{
Create(data *%s) error
Get(id int) (%s,error)
Find(condition interface{},page,limit int) ([]%s,error)
Delete(id int) error
DeleteUnScope(id int) error
}
`
	info.WriteString(fmt.Sprintf(interfaceContent, interfaceName, name, name, name))
	info.WriteString("\n")
	// 2

	modelServiceName := strings.ToLower(name[0:1]) + name[1:] + "ModelService"
	modelService := `
type %s struct{
db *gorm.DB
}
`
	info.WriteString(fmt.Sprintf(modelService, modelServiceName))
	info.WriteString("\n")
	newModelService := `
func New%sModelService(db *gorm.DB) %s {
return %s{db:db}
}
`
	info.WriteString(fmt.Sprintf(newModelService, name, interfaceName, modelServiceName))
	info.WriteString("\n")

	// 3
	create := `
func (s %s) Create(data *%s) error{
err:=s.db.Create(data).Error
if err != nil{
return err
}
return nil
} 
`
	info.WriteString(fmt.Sprintf(create, modelServiceName, name))
	info.WriteString("\n")

	//func (s userModelService) Get(id int) (User, error) {
	//	var u User
	//	err := s.db.Where(id).Find(&u).Limit(1).Error
	//	if err != nil {
	//		return User{}, err
	//	}
	//	return u, nil
	//}

	get := `
func (s %s) Get(id int) (%s,error){
var u %s
err:=s.db.Where(id).Find(&u).Limit(1).Error
if err != nil{
return %s{},err
}
return u,nil
} 
`
	info.WriteString(fmt.Sprintf(get, modelServiceName, name, name, name))
	info.WriteString("\n")

	find := `
func (s %s) Find(condition interface{},page,limit int) ([]%s,error){
var list []%s
err:=s.db.Where(condition).Find(&list).Offset(limit * (page - 1)).Limit(limit).Error
if err != nil{
return nil,err
}
return list,nil
} 
`
	info.WriteString(fmt.Sprintf(find, modelServiceName, name, name))
	info.WriteString("\n")

	del := `
func (s %s) Delete(id int)  error {
err:=s.db.Delete(id).Error
if err != nil{
return  err
}
return  nil
} 
`
	info.WriteString(fmt.Sprintf(del, modelServiceName))
	info.WriteString("\n")

	delUnScope := `
func (s %s) DeleteUnScope(id int)  error {
err:=s.db.Unscoped().Delete(id).Error
if err != nil{
return  err
}
return  nil
} 
`
	info.WriteString(fmt.Sprintf(delUnScope, modelServiceName))
	info.WriteString("\n")
	return info.String()
}

// DealColumn judge column type and generate tag info
func (s service) dealColumn(t *tableInfo) {
	for i := 0; i < len(t.column); i++ {
		var f bool
		if s.Conf.IsGenJsonTag {
			//生成 json tag
			f = true
			t.column[i].Tag = "`json:\"" + JsonTag(s.Conf.JsonTagType, t.column[i].ColumnName) + "\""
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
		t.column[i].ColumnName = UpperCamel(t.column[i].ColumnName)
		t.column[i].ColumnType = s.dealType(s.Conf, t.column[i].DataType, t.column[i].ColumnType)
	}
}

func (s service) dealStructContent(t tableInfo) string {
	var info strings.Builder
	// struct name
	structName := UpperCamel(t.TableName)

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

func (s service) dealType(c Config, typeSimple, typeDetail string) string {
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