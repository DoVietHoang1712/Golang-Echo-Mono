package store

import (
	"golang-sql/model"
	"gorm.io/gorm"
)

type UserStore struct {
	db *gorm.DB
}

func NewUserStore(db *gorm.DB) *UserStore {
	return &UserStore{
		db: db,
	}
}

func (us *UserStore) GetByID(id uint) (*model.User, error) {
	var u model.User
	if err := us.db.First(&u, id).Error; err != nil {
		return nil, err
	}
	return &u, nil
}

func (us *UserStore) GetByEmail(email string) (*model.User, error) {
	var u model.User
	if err := us.db.Where(&model.User{Email: email}).First(&u).Error; err != nil {
		return nil, err
	}
	return &u, nil
}

func (us *UserStore) GetByUsername(username string) (*model.User, error) {
	var u model.User
	if err := us.db.Where(&model.User{Username: username}).Preload("Followers").First(&u).Error; err != nil {
		return nil, err
	}
	return &u, nil
}

func (us *UserStore) Create(user *model.User) error {
	return us.db.Create(user).Error
}

func (us *UserStore) Update(user *model.User) error {
	return us.db.Model(user).Updates(user).Error
}
