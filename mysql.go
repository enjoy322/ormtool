package ormtool

import (
	"database/sql"
	"fmt"
	"log"
	"strings"
)

type serviceInterface interface {
	genStruct() (fileData FileInfo, data []StructInfo)

	listTables() []tableInfo

	listColumns() map[string][]column

	dealColumn(t *tableInfo)

	dealStructContent(t tableInfo) string

	dealType(c Config, typeSimple, typeDetail string) string

	getCreateSQL(tableName string) string

	gen()
}

func GenerateMySQL(c Config) {
	newService(c).gen()
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

func newService(c Config) serviceInterface {
	conn, err := mysqlConn(c.ConnStr)
	if err != nil {
		log.Fatalln(err)
	}

	return service{DB: conn, dbName: c.Database, Conf: c}
}

func (s service) gen() {
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
	//ColumnKey string
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

		info.FileName = lowerCamel(table.TableName)

		info.Name = UpperCamel(table.TableName)

		// table comment
		// add if table comment exists
		if table.TableComment != "" || s.Conf.IsGenCreateSQL {
			info.Note = "// " + info.Name + "\t" + table.TableComment + "\n"
		}

		if s.Conf.IsGenFunction {
			//	simple function
			if s.Conf.IsGenFunctionWithCache {
				info.ImportInfo = []string{"context", "encoding/json", "strconv", "github.com/go-redis/redis/v8",
					"gorm.io/gorm", "time"}
				info.Function = s.genFunctionWithCache(info.Name)
			} else {
				info.ImportInfo = []string{"gorm.io/gorm"}
				if s.Conf.IsGenFunctionWithCache {
					info.ImportInfo = append(info.ImportInfo, "context")
				}

				info.Function = s.genFunction(info.Name)
			}
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
Create(tx *gorm.DB,data *%s) error
Get(tx *gorm.DB,id int) (%s,error)
Find(tx *gorm.DB,page,limit int) ([]%s,int64,error)
DeleteByID(tx *gorm.DB,id int) error
}
`
	info.WriteString(fmt.Sprintf(interfaceContent, interfaceName, name, name, name))
	info.WriteString("\n")
	// 2

	modelServiceName := strings.ToLower(name[0:1]) + name[1:] + "ModelService"
	modelService := `
type %s struct{
}
`
	info.WriteString(fmt.Sprintf(modelService, modelServiceName))
	info.WriteString("\n")
	newModelService := `
func New%sModelService() %s {
return %s{}
}
`
	info.WriteString(fmt.Sprintf(newModelService, name, interfaceName, modelServiceName))
	info.WriteString("\n")

	// 3
	create := `
func (s %s) Create(tx *gorm.DB,data *%s) error{
err:=tx.Create(data).Error
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
func (s %s) Get(tx *gorm.DB,id int) (%s,error){
var data %s
err:=tx.Where(id).First(&data).Error
if err != nil && err==gorm.ErrRecordNotFound {
return %s{},err
}
return data,nil
} 
`
	info.WriteString(fmt.Sprintf(get, modelServiceName, name, name, name))
	info.WriteString("\n")

	find := `
func (s %s) Find(tx *gorm.DB,page,limit int) ([]%s,int64,error){
var list []%s
err:=tx.Offset(limit * (page - 1)).Limit(limit).Find(&list).Error
if err != nil{
return nil,0,err
}
var count int64
err = tx.Count(&count).Error
if err != nil {
	return nil,0,err
}
return list,count,nil
} 
`
	info.WriteString(fmt.Sprintf(find, modelServiceName, name, name))
	info.WriteString("\n")

	del := `
func (s %s) DeleteByID(tx *gorm.DB,id int)  error {
err:=tx.Where(id).Delete(&%s{}).Error
if err != nil{
return  err
}
return  nil
} 
`
	info.WriteString(fmt.Sprintf(del, modelServiceName, name))
	info.WriteString("\n")

	return info.String()
}

func (s service) genFunctionWithCache(name string) string {
	var info strings.Builder
	camelName := strings.ToLower(name[0:1]) + name[1:]
	info.WriteString("// function\n")
	cacheName := fmt.Sprintf("%sCache", camelName)
	invalidCacheName := fmt.Sprintf("%sInvalidCache", camelName)
	info.WriteString(fmt.Sprintf("var %s=\"cache%s:\"\n", cacheName, name))
	info.WriteString(fmt.Sprintf("var %s=\"invalidCache%s:\"\n", invalidCacheName, name))

	interfaceName := name + "ModelInterface"
	interfaceContent := `
type %s interface{
Create(tx *gorm.DB,data *%s) error
Get(tx *gorm.DB,id int) (%s,error)
Find(tx *gorm.DB,page,limit int) ([]%s,int64,error)
DeleteByID(tx *gorm.DB,id int) error
DeleteCache(id int)
}
`
	info.WriteString(fmt.Sprintf(interfaceContent, interfaceName, name, name, name))
	info.WriteString("\n")
	// 2

	modelServiceName := camelName + "ModelService"
	modelService := `
type %s struct{
rdb *redis.Client
}
`
	info.WriteString(fmt.Sprintf(modelService, modelServiceName))
	info.WriteString("\n")
	newModelService := `
func New%sModelService(redisDB *redis.Client) %s {
return %s{rdb:redisDB}
}
`
	info.WriteString(fmt.Sprintf(newModelService, name, interfaceName, modelServiceName))
	info.WriteString("\n")

	// 3
	create := `
func (s %s) Create(tx *gorm.DB,data *%s) error{
err:=tx.Create(data).Error
if err != nil{
return err
}
marshal, err := json.Marshal(data)
	if err != nil {
		return err
	}
s.rdb.Set(context.Background(),%s+strconv.Itoa(data.Id),string(marshal),time.Hour*48)
return nil
} 
`
	info.WriteString(fmt.Sprintf(create, modelServiceName, name, cacheName))
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
func (s %s) Get(tx *gorm.DB,id int) (%s,error){
invalidKey := %s + strconv.Itoa(id)
if s.rdb.Exists(context.Background(), invalidKey).Val() > 0 {	
	return %s{}, nil
}
var data %s
key := %s + strconv.Itoa(id)
if s.rdb.Exists(context.Background(), key).Val() > 0 {
	bytes, err := s.rdb.Get(context.Background(), key).Bytes()
	if err != nil {
		return %s{}, err
	}
	err = json.Unmarshal(bytes, &data)
	if err != nil {
		return %s{}, err
	}
	return data, nil
}
err:=tx.Where(id).First(&data).Error
if err != nil  && err!=gorm.ErrRecordNotFound{
return %s{},err
}
if err==gorm.ErrRecordNotFound{
return %s{},nil
}
if err != gorm.ErrRecordNotFound {
	//	exist
	marshal, err := json.Marshal(data)
	if err != nil {
		return data,err
	}
	s.rdb.Set(context.Background(),%s+strconv.Itoa(data.Id),string(marshal),time.Hour*48)
	return data,nil
}
s.rdb.Set(context.Background(),invalidKey,"",time.Minute*2)
return data,nil
} 
`
	info.WriteString(fmt.Sprintf(get, modelServiceName, name, invalidCacheName, name, name, cacheName, name, name, name, name, cacheName))
	info.WriteString("\n")

	find := `
func (s %s) Find(tx *gorm.DB,page,limit int) ([]%s,int64,error){
var list []%s
var count int64
err := tx.Count(&count).Error
if err != nil {
	return nil,0,err
}
err=tx.Offset(limit * (page - 1)).Limit(limit).Find(&list).Error
if err != nil{
return nil,0,err
}

return list,count,nil
} 
`
	info.WriteString(fmt.Sprintf(find, modelServiceName, name, name))
	info.WriteString("\n")

	del := `
func (s %s) DeleteByID(tx *gorm.DB,id int)  error {
err:=tx.Where(id).Delete(&%s{}).Error
if err != nil{
return  err
}
s.rdb.Del(context.Background(),%s+strconv.Itoa(id))
return  nil
} 
`
	info.WriteString(fmt.Sprintf(del, modelServiceName, name, cacheName))
	info.WriteString("\n")

	delCache := `
func (s %s) DeleteCache(id int){
s.rdb.Del(context.Background(),%s+strconv.Itoa(id))
}
`
	info.WriteString(fmt.Sprintf(delCache, modelServiceName, cacheName))
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

// GetCreateSQL sql of creating table in the database
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
   FROM information_schema.COLUMNS WHERE table_schema = ? order by ORDINAL_POSITION`
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
