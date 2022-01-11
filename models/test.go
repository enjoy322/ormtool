package models

// User
/*CREATE TABLE `user` (
  `id` int NOT NULL AUTO_INCREMENT,
  `create_time` int DEFAULT NULL COMMENT '创建时间',
  `user_name` varchar(2550) DEFAULT 'qw' COMMENT '用户名',
  `u2` varchar(127) DEFAULT NULL,
  `u3` varchar(255) DEFAULT NULL,
  `u4` varchar(256) DEFAULT NULL,
  `u5` int DEFAULT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `user_u3_uindex` (`u3`),
  UNIQUE KEY `user_u5_uindex` (`u5`),
  KEY `user_u2_index` (`u2`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb3 COMMENT='用户'*/
type User struct {
	CreateTime int    `json:"create_time" ` // 创建时间
	Id         int    `json:"id" `
	U2         string `json:"u2" `
	U3         string `json:"u3" `
	U4         string `json:"u4" `
	U5         int    `json:"u5" `
	UserName   string `json:"user_name" ` // 用户名
}

func (*User) TableName() string {
	return "user"
}

var UserCol = struct {
	CreateTime string
	Id         string
	U2         string
	U3         string
	U4         string
	U5         string
	UserName   string
}{
	CreateTime: "create_time",
	Id:         "id",
	U2:         "u2",
	U3:         "u3",
	U4:         "u4",
	U5:         "u5",
	UserName:   "user_name",
}

// UserCard
/*CREATE TABLE `user_card` (
  `id` int NOT NULL AUTO_INCREMENT,
  `addr` varchar(128) DEFAULT NULL,
  `is_banned` tinyint(1) DEFAULT '0',
  `a3` char(1) DEFAULT NULL,
  `a4` char(12) DEFAULT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `user_card_id_uindex` (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci*/
type UserCard struct {
	A3       string `json:"a3" `
	A4       string `json:"a4" `
	Addr     string `json:"addr" `
	Id       int    `json:"id" `
	IsBanned bool   `json:"is_banned" `
}

func (*UserCard) TableName() string {
	return "user_card"
}

var UserCardCol = struct {
	A3       string
	A4       string
	Addr     string
	Id       string
	IsBanned string
}{
	A3:       "a3",
	A4:       "a4",
	Addr:     "addr",
	Id:       "id",
	IsBanned: "is_banned",
}

// UserCardDetail
/*CREATE TABLE `user_card_detail` (
  `user_card_id` int DEFAULT NULL,
  KEY `user_card_detail_user_card_id_fk` (`user_card_id`),
  CONSTRAINT `user_card_detail_user_card_id_fk` FOREIGN KEY (`user_card_id`) REFERENCES `user_card` (`id`) ON DELETE CASCADE ON UPDATE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci*/
type UserCardDetail struct {
	UserCardId int `json:"user_card_id" `
}

func (*UserCardDetail) TableName() string {
	return "user_card_detail"
}

var UserCardDetailCol = struct {
	UserCardId string
}{
	UserCardId: "user_card_id",
}
