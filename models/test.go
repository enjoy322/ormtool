package models

type User struct {
	Id         int    `json:"id" db:"int"`
	CreateTime int    `json:"create_time" db:"int"`        // 创建时间
	UserName   string `json:"user_name" db:"varchar(255)"` // 用户名
}

func (User) TableName() string {
	return "user"
}

type Book struct {
	Id       int     `json:"id" db:"int"`
	BookName int     `json:"book_name" db:"int"`
	Price    float64 `json:"price" db:"double"`
	IsBanned bool    `json:"is_banned" db:"tinyint(1)"`
}

func (Book) TableName() string {
	return "book"
}
