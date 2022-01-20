package models

// User
/*CREATE TABLE `user` (
  `id` int unsigned NOT NULL AUTO_INCREMENT COMMENT '主键',
  `mobile` char(12) DEFAULT NULL COMMENT '手机号',
  `pwd` varchar(256) DEFAULT NULL COMMENT '密码',
  PRIMARY KEY (`id`),
  UNIQUE KEY `user_id_uindex` (`id`),
  UNIQUE KEY `user_mobile_uindex` (`mobile`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='用户表'*/
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
