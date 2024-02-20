// Code generated by ormtool. DO NOT EDIT.
package model

type Question struct {
	Id            int    `json:"id"`
	Title         string `json:"title"`          // 标题
	UserId        int    `json:"user_id"`        // 用户ID
	Content       string `json:"content"`        // 内容
	ReceiveAnswer int    `json:"receive_answer"` // 接收答案1 接收 2 不接收
	Status        int    `json:"status"`         // 状态
	CreatedAt     int    `json:"created_at"`
	UpdatedAt     int    `json:"updated_at"`
	DeletedAt     int    `json:"deleted_at"`
}

var QuestionCol = struct {
	Id            string
	Title         string
	UserId        string
	Content       string
	ReceiveAnswer string
	Status        string
	CreatedAt     string
	UpdatedAt     string
	DeletedAt     string
}{
	Id:            "id",
	Title:         "title",
	UserId:        "user_id",
	Content:       "content",
	ReceiveAnswer: "receive_answer",
	Status:        "status",
	CreatedAt:     "created_at",
	UpdatedAt:     "updated_at",
	DeletedAt:     "deleted_at",
}
