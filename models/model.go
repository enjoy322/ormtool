package models

type TableName struct {
	Id int `json:"id" ` // 主键
}

func (*TableName) TableName() string {
	return "table_name"
}

var TableNameCol = struct {
	Id string
}{
	Id: "id",
}

// User	用户表
type User struct {
	Id     uint32 `json:"id" `     // 主键
	Mobile string `json:"mobile" ` // 手机号
	Pwd    string `json:"pwd" `    // 密码
}

func (*User) TableName() string {
	return "user"
}

var UserCol = struct {
	Id     string
	Mobile string
	Pwd    string
}{
	Id:     "id",
	Mobile: "mobile",
	Pwd:    "pwd",
}
