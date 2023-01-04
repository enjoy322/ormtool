package model

type UserCard struct {
	A3       string `json:"a3"`
	A4       string `json:"a4"`
	Addr     string `json:"addr"`
	Id       int    `json:"id"`
	IsBanned bool   `json:"is_banned"`
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

// function

type UserCardModelInterface interface {
	Create(data *UserCard) error
	Get(id int) (UserCard, error)
	Find(condition *gorm.DB, page, limit int) ([]UserCard, int64, error)
	Delete(id int) error
	DeleteUnScope(id int) error
}

type userCardModelService struct {
	db *gorm.DB
}

func NewUserCardModelService(db *gorm.DB) UserCardModelInterface {
	return userCardModelService{db: db}
}

func (s userCardModelService) Create(data *UserCard) error {
	err := s.db.Create(data).Error
	if err != nil {
		return err
	}
	return nil
}

func (s userCardModelService) Get(id int) (UserCard, error) {
	var u UserCard
	err := s.db.Where(id).Find(&u).Limit(1).Error
	if err != nil {
		return UserCard{}, err
	}
	return u, nil
}

func (s userCardModelService) Find(condition *gorm.DB, page, limit int) ([]UserCard, int64, error) {
	var list []UserCard
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

func (s userCardModelService) Delete(id int) error {
	err := s.db.Delete(id).Error
	if err != nil {
		return err
	}
	return nil
}

func (s userCardModelService) DeleteUnScope(id int) error {
	err := s.db.Unscoped().Delete(id).Error
	if err != nil {
		return err
	}
	return nil
}
