package ormtool

// Config config information
type Config struct {
	// connect mysql, [user]:[password]@tcp([host]:[port])/[database]?parseTime=true
	ConnStr string
	// database name
	Database string
	// file, relative path
	SavePath string
	// json tag
	IsGenJsonTag bool
	// Generate one file or files by table
	IsGenInOneFile bool
	// Generate simple database field information like: "int unsigned not null"
	// value 1:not generate; 2：simple info
	GenDBInfoType int
	// json tag type. The necessary conditions：IsGenJsonTag:true.
	// 1.UserName 2.userName 3.user_name 4.user-name
	JsonTagType int
	// sql of creating table in the database
	IsGenCreateSQL bool
	// simple crud function
	IsGenFunction bool
	// cache simple model info to redis, Ps. IsGenCreateSQL = true
	IsGenFunctionWithCache bool
	// method TableName
	IsGenTableName bool
	// custom type relationships will be preferred.
	// the key is the database type, The value is the golang type
	CustomType map[string]string
}

type StructInfo struct {
	// package信息
	ImportInfo    []string
	Name          string
	FileName      string //model file name
	TableName     string
	Note          string // description
	CreateSQL     string // create table sql
	StructContent string
	Function      string // simple crud function
}

type FileInfo struct {
	PackageName string
	FileDir     string
	FileName    string
}
