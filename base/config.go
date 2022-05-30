package base

// Config 配置文件
type Config struct {
	// 生成文件名，默认[数据库].go,使用相对路径
	SavePath string
	// 生成json tag
	IsGenJsonTag bool
	// 生成在同一个文件中
	IsGenInOneFile bool
	// 生成数据库字段信息  1.不生产 2.普通字段信息
	GenDBInfoType int
	// jsonTag类型 1.UserName 2.userName 3.user_name 4.user-name
	JsonTagType int
	//	是否生成建表语句
	IsGenCreateSQL bool
	//	自定义数据库和Go类型对应关系
	CustomType map[string]string
}

// MysqlConfig mysql配置
type MysqlConfig struct {
	User     string
	Password string
	Host     string
	Port     string
	// 数据库名称
	Database string
}

type StructInfo struct {
	Name          string
	TableName     string
	Note          string // descript
	CreateSQL     string // create table sql
	StructContent string
}

type FileInfo struct {
	PackageName string
	FileDir     string
	FileName    string
}
