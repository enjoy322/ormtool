package models

import "encoding/json"

// User	用户表
/*CREATE TABLE `user` (
  `id` int NOT NULL AUTO_INCREMENT,
  `user_name` varchar(64) DEFAULT NULL COMMENT '用户名',
  `avatar` json DEFAULT NULL COMMENT '头像',
  `gender` tinyint(1) DEFAULT '1' COMMENT '性别',
  PRIMARY KEY (`id`),
  UNIQUE KEY `user_id_uindex` (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='用户表'*/
type User struct {
	Id       int             `json:"id" `
	UserName string          `json:"user_name" ` // 用户名
	Avatar   json.RawMessage `json:"avatar" `    // 头像
	Gender   bool            `json:"gender" `    // 性别
}

func (*User) TableName() string {
	return "user"
}

var UserCol = struct {
	Id       string
	UserName string
	Avatar   string
	Gender   string
}{
	Id:       "id",
	UserName: "user_name",
	Avatar:   "avatar",
	Gender:   "gender",
}
