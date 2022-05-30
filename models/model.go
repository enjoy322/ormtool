package models

// User	用户表
/*CREATE TABLE `user` (
  `id` int unsigned NOT NULL COMMENT '主键',
  `created_at` int unsigned DEFAULT NULL COMMENT '创建时间',
  `updated_at` int unsigned DEFAULT NULL COMMENT '更新时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `user_id_uindex` (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='用户表'*/
type User struct {
	Id        uint32 `json:"id" db:"int unsigned not null"` // 主键
	CreatedAt uint32 `json:"created_at" db:"int unsigned"`  // 创建时间
	UpdatedAt uint32 `json:"updated_at" db:"int unsigned"`  // 更新时间
}

func (*User) TableName() string {
	return "user"
}

var UserCol = struct {
	Id        string
	CreatedAt string
	UpdatedAt string
}{
	Id:        "id",
	CreatedAt: "created_at",
	UpdatedAt: "updated_at",
}
