package base

// Config 配置文件
type Config struct {
	//文件名，默认[数据库].go,使用相对路径
	SavePath string
	//Table   []string todo
	IsGenJsonTag bool
	//生成在同一个文件中
	IsGenInOneFile bool
	//生成数据库字段信息  1.不生产 2.普通字段信息 3.gorm 4.xorm
	GenDBInfoType int
	//jsonTag类型 1.UserName 2.user_name 3.userName 4.user-name
	JsonTagType int
	//	是否生成建表语句
	IsGenCreateSQL bool
	//	CustomType
	CustomType map[string]string
}

// MysqlConfig mysql配置
type MysqlConfig struct {
	User     string
	Password string
	Host     string
	Port     string
	Database string
}
