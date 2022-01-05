package base

const (
	//UserName
	CodeJsonTag1 = 1 + iota
	//user_name
	CodeJsonTag2
	//userName
	CodeJsonTag3
	//user-name
	CodeJsonTag4
)

const (
	CodeDBMySQL = 1 + iota
	CodeDBMSSQL
)

const (
	CodeDBInfoSimple = 1 + iota
	CodeDBInfoGorm
	CodeDBInfoXorm
)

type Config struct {
	DataBaseType int
	MySQL        MysqlConfig
	//文件名，默认[数据库].go,使用相对路径
	SavePath string
	//Table   []string todo
	IsGenJsonTag bool
	//生成在同一个文件中
	IsGenInOneFile bool
	//是否生成数据库字段信息
	GenDBInfoType int
	JsonTagType   int
}

type MysqlConfig struct {
	User     string
	Password string
	Host     string
	Port     string
	Database string
}
