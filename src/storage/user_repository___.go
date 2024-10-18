package storage

import (
	webmodels "github.com/pa-pe/luca-control/src/web/models"
)

type UserRepository interface {
	CreateUser(user *webmodels.WebUser) error
	GetUserByID(id uint) (*webmodels.WebUser, error)
	ListUsers() ([]*webmodels.WebUser, error)
}
