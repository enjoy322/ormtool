package models

// User	用户表
type User struct {
	CreateTime int    `json:"create_time" ` // 创建时间
	Id         int    `json:"id" `
	Temp       int    `json:"temp" `
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
	Temp       string
	U2         string
	U3         string
	U4         string
	U5         string
	UserName   string
}{
	CreateTime: "create_time",
	Id:         "id",
	Temp:       "temp",
	U2:         "u2",
	U3:         "u3",
	U4:         "u4",
	U5:         "u5",
	UserName:   "user_name",
}

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
