package model

type UserCardDetail struct {
	UserCardId int `json:"user_card_id"`
}

func (*UserCardDetail) TableName() string {
	return "user_card_detail"
}

var UserCardDetailCol = struct {
	UserCardId string
}{
	UserCardId: "user_card_id",
}

// function

type UserCardDetailModelInterface interface {
	Create(data *UserCardDetail) error
	Get(id int) (UserCardDetail, error)
	Find(condition *gorm.DB, page, limit int) ([]UserCardDetail, int64, error)
	Delete(id int) error
	DeleteUnScope(id int) error
}

type userCardDetailModelService struct {
	db *gorm.DB
}

func NewUserCardDetailModelService(db *gorm.DB) UserCardDetailModelInterface {
	return userCardDetailModelService{db: db}
}

func (s userCardDetailModelService) Create(data *UserCardDetail) error {
	err := s.db.Create(data).Error
	if err != nil {
		return err
	}
	return nil
}

func (s userCardDetailModelService) Get(id int) (UserCardDetail, error) {
	var u UserCardDetail
	err := s.db.Where(id).Find(&u).Limit(1).Error
	if err != nil {
		return UserCardDetail{}, err
	}
	return u, nil
}

func (s userCardDetailModelService) Find(condition *gorm.DB, page, limit int) ([]UserCardDetail, int64, error) {
	var list []UserCardDetail
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

func (s userCardDetailModelService) Delete(id int) error {
	err := s.db.Delete(id).Error
	if err != nil {
		return err
	}
	return nil
}

func (s userCardDetailModelService) DeleteUnScope(id int) error {
	err := s.db.Unscoped().Delete(id).Error
	if err != nil {
		return err
	}
	return nil
}
