package models

import "encoding/json"

// A	a
/*CREATE TABLE `a` (
  `a1` int DEFAULT NULL COMMENT 'a11111'
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='a'*/
type A struct {
	A1 int `json:"a1"` // a11111
}

func (*A) TableName() string {
	return "a"
}

var ACol = struct {
	A1 string
}{
	A1: "a1",
}

// C	c
/*CREATE TABLE `c` (
  `c1` int DEFAULT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='c'*/
type C struct {
	C1 int `json:"c1"`
}

func (*C) TableName() string {
	return "c"
}

var CCol = struct {
	C1 string
}{
	C1: "c1",
}

// F
/*CREATE TABLE `f` (
  `f1` int DEFAULT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci*/
type F struct {
	F1 int `json:"f1"`
}

func (*F) TableName() string {
	return "f"
}

var FCol = struct {
	F1 string
}{
	F1: "f1",
}

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
	Id       int             `json:"id"`
	UserName string          `json:"user_name"` // 用户名
	Avatar   json.RawMessage `json:"avatar"`    // 头像
	Gender   bool            `json:"gender"`    // 性别
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

// W	w
/*CREATE TABLE `w` (
  `w1` int DEFAULT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='w'*/
type W struct {
	W1 int `json:"w1"`
}

func (*W) TableName() string {
	return "w"
}

var WCol = struct {
	W1 string
}{
	W1: "w1",
}

// X	x
/*CREATE TABLE `x` (
  `x1` int DEFAULT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='x'*/
type X struct {
	X1 int `json:"x1"`
}

func (*X) TableName() string {
	return "x"
}

var XCol = struct {
	X1 string
}{
	X1: "x1",
}
