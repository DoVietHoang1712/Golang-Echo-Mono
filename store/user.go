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

func (us *UserStore) IsFollower(userID, followerID uint) (bool, error) {
	var f model.Follow
	if err := us.db.Where("following_id = ? AND follower_id = ?", userID, followerID).Find(&f).Error; err != nil {
		return false, err
	}
	return true, nil
}

func (us *UserStore) AddFollower(user *model.User, followerID uint) error {
	return us.db.Model(user).Association("Follower").Append(&model.Follow{
		FollowerID:  followerID,
		FollowingID: user.ID,
	})
}

func (us *UserStore) RemoveFollower(user *model.User, followerID uint) error {
	f := model.Follow{
		FollowerID:  followerID,
		FollowingID: user.ID,
	}
	if err := us.db.Model(user).Association("Follower").Find(&f); err != nil {
		return err
	}
	if err := us.db.Delete(f).Error; err != nil {
		return err
	}
	return nil
}
