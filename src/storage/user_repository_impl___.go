package storage

import (
	"github.com/pa-pe/luca-control/src/storage/model"
	"gorm.io/gorm"
)

type UserRepositoryImpl struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return &UserRepositoryImpl{db: db}
}

func (r *UserRepositoryImpl) CreateUser(user *model.WebUser) error {
	return r.db.Create(user).Error
}

func (r *UserRepositoryImpl) GetUserByID(id uint) (*model.WebUser, error) {
	var user model.WebUser
	if err := r.db.First(&user, id).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *UserRepositoryImpl) ListUsers() ([]*model.WebUser, error) {
	var users []*model.WebUser
	if err := r.db.Find(&users).Error; err != nil {
		return nil, err
	}
	return users, nil
}
