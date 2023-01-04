package model

// User	用户表
type User struct {
	CreateTime int    `json:"create_time"` // 创建时间
	Id         int    `json:"id"`
	Temp       int    `json:"temp"`
	U2         string `json:"u2"`
	U3         string `json:"u3"`
	U4         string `json:"u4"`
	U5         int    `json:"u5"`
	Un         uint32 `json:"un"`
	Un2        uint32 `json:"un2"`
	UserName   string `json:"user_name"` // 用户名
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
	Un         string
	Un2        string
	UserName   string
}{
	CreateTime: "create_time",
	Id:         "id",
	Temp:       "temp",
	U2:         "u2",
	U3:         "u3",
	U4:         "u4",
	U5:         "u5",
	Un:         "un",
	Un2:        "un2",
	UserName:   "user_name",
}

// function

type UserModelInterface interface {
	Create(data *User) error
	Get(id int) (User, error)
	Find(condition *gorm.DB, page, limit int) ([]User, int64, error)
	Delete(id int) error
	DeleteUnScope(id int) error
}

type userModelService struct {
	db *gorm.DB
}

func NewUserModelService(db *gorm.DB) UserModelInterface {
	return userModelService{db: db}
}

func (s userModelService) Create(data *User) error {
	err := s.db.Create(data).Error
	if err != nil {
		return err
	}
	return nil
}

func (s userModelService) Get(id int) (User, error) {
	var u User
	err := s.db.Where(id).Find(&u).Limit(1).Error
	if err != nil {
		return User{}, err
	}
	return u, nil
}

func (s userModelService) Find(condition *gorm.DB, page, limit int) ([]User, int64, error) {
	var list []User
	err := condition.Find(&list).Offset(limit * (page - 1)).Limit(limit).Error
	if err != nil {
		return nil, 0, err
	}
	var count int64
	err = condition.Count(&count).Error
	if err != nil {
		return nil, 0, err
	}
	return list, count, nil
}

func (s userModelService) Delete(id int) error {
	err := s.db.Delete(id).Error
	if err != nil {
		return err
	}
	return nil
}

func (s userModelService) DeleteUnScope(id int) error {
	err := s.db.Unscoped().Delete(id).Error
	if err != nil {
		return err
	}
	return nil
}
