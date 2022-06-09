package base

// Config config information
type Config struct {
	// file relative path
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
	// sql of creating table in database
	IsGenCreateSQL bool
	// custom type relationships will be preferred
	// the key is the database type, the value is the golang type
	CustomType map[string]string
}

// MysqlConfig mysql config information
type MysqlConfig struct {
	User     string
	Password string
	Host     string
	Port     string
	// database name
	Database string
}

type StructInfo struct {
	Name          string
	TableName     string
	Note          string // description
	CreateSQL     string // create table sql
	StructContent string
}

type FileInfo struct {
	PackageName string
	FileDir     string
	FileName    string
}
