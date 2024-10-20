package storage

import (
	webmodels "github.com/pa-pe/luca-control/src/storage/model"
)

type UserRepository interface {
	CreateUser(user *webmodels.WebUser) error
	GetUserByID(id uint) (*webmodels.WebUser, error)
	ListUsers() ([]*webmodels.WebUser, error)
}
